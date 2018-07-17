package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"strings"
	"time"
)

var (
	hosts        = flag.String("hosts", "localhost", "Comma-separated hosts and IP addresses to generate the certificate for")
	expiration   = flag.Duration("expiration", 365*24*time.Hour, "Duration that certificate is valid for")
	size         = flag.Int("bits", 2048, "Size of key to generate.")
	organization = flag.String("org", "Acme Testing Corp", "Organization name")
	certOut      = flag.String("cert", "cert.pem", "File to write the certificate to")
	keyOut       = flag.String("key", "key.pem", "File to write the key to")
)

func main() {
	flag.Parse()

	pk, err := rsa.GenerateKey(rand.Reader, *size)
	if err != nil {
		log.Fatalf("Could not generate private key: %v", err)
	}

	if *hosts == "" {
		fmt.Fprint(os.Stderr, "Missing required hosts parameters")
		flag.PrintDefaults()
		os.Exit(-2)
	}
	hostList := strings.Split(*hosts, ",")
	var names []string
	var ips []net.IP
	for _, h := range hostList {
		if ip := net.ParseIP(h); ip != nil {
			ips = append(ips, ip)
			continue
		}
		names = append(names, h)
	}

	template := x509.Certificate{
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,

		Subject: pkix.Name{
			Organization: []string{*organization},
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(*expiration),
		SerialNumber: serialNumber(),

		IPAddresses: ips,
		DNSNames:    names,
	}

	cert, err := x509.CreateCertificate(rand.Reader, &template, &template, &pk.PublicKey, pk)
	if err != nil {
		log.Fatalf("Could not create certificate: %v", err)
	}

	f, err := os.Create(*certOut)
	if err != nil {
		log.Fatalf("Could not write to %s: %s", *certOut, err)
	}
	pem.Encode(f, &pem.Block{Type: "CERTIFICATE", Bytes: cert})
	f.Close()
	fmt.Printf("Certificate written to %s\n", *certOut)

	f, err = os.OpenFile("key.pem", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Print("Failed to open key.pem for writing:", err)
		return
	}
	pem.Encode(f, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
	f.Close()
	fmt.Printf("Key written to %s\n", *keyOut)
}

// serialNumber returns a random number between 0 and 2^160 that can be used as
// the serial number for a certificate
func serialNumber() *big.Int {
	var max big.Int
	sn, err := rand.Int(rand.Reader, max.Exp(big.NewInt(2), big.NewInt(160), nil))
	if err != nil {
		log.Fatalf("Could not generate serial number: %v; returning 0", err)
	}

	return sn
}
