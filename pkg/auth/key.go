package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func LoadPrivateKeyFromFile(filepath string) (*rsa.PrivateKey, error) {
	// Read the bytes of the PEM file, e.g. id_rsa
	pemData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	// Use the PEM decoder and parse the private key
	pemBlock, _ := pem.Decode(pemData)
	priv, err := x509.ParsePKCS1PrivateKey(pemBlock.Bytes)

	// Public key can be obtained through priv.PublicKey
	return priv, err
}

func LoadPrivateKeyOrDie(filepath string) *rsa.PrivateKey {
	p, err := LoadPrivateKeyFromFile(filepath)
	if err != nil {
		panic("Cannot load PEM data")
	}
	return p
}
