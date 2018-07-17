# Keygen

Simple helper tool for creating test, self-signed certificates. By default, it will write
two files: cert.pem and key.pem.

## Installation

    go get github.com/rselbach/keygen/cmd/keygen

## Usage

    -bits int
    	Size of RSA key to generate. (default 2048)
    -cert string
    	File to write the certificate to (default "cert.pem")
    -expiration duration
    	Duration that certificate is valid for (default 1 year)
    -hosts string
    	Comma-separated hosts and IP addresses to generate the certificate for (default "localhost")
    -key string
    	File to write the key to (default "key.pem")
    -org string
    	Organization name (default "Acme Co.")


## Examples

Default usage, creating a certificate that last a whole year for localhost:

    $ keygen
    Certificate written to cert.pem
    Key written to key.pem

Generate a certificate for multiple hosts:

    $ keygen -hosts example.com,10.213.200.7
    Certificate written to cert.pem
    Key written to key.pem


A certificate that expires in 10 minutes:

    $ keygen -hosts example.com -expiration 10m
    Certificate written to cert.pem
    Key written to key.pem

Specifying file names:

    $ keygen -hosts example.com -cert example.crt -key example.key
    Certificate written to example.crt
    Key written to example.key

