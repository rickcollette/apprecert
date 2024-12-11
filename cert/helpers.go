package cert

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"

	"github.com/appremon/apprecert/config"
)

// GenerateKey generates a new private key (RSA or ECDSA).
func GenerateKey(ecdsaKey bool, rootCA bool) (crypto.PrivateKey, error) {
	if ecdsaKey {
		return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	}
	if rootCA {
		return rsa.GenerateKey(rand.Reader, 3072)
	}
	return rsa.GenerateKey(rand.Reader, 2048)
}

// RandomSerialNumber generates a random serial number for a certificate.
func randomSerialNumber() *big.Int {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		panic("failed to generate serial number: " + err.Error())
	}
	return serialNumber
}

// saveCertificate saves the generated certificate and private key to disk.
func saveCertificate(cfg *config.Config, certBytes []byte, hosts []string) error {
	certPath := filepath.Join(cfg.CAROOT, fmt.Sprintf("%s-cert.pem", hosts[0]))
	keyPath := filepath.Join(cfg.CAROOT, fmt.Sprintf("%s-key.pem", hosts[0]))

	// Save certificate
	certFile, err := os.Create(certPath)
	if err != nil {
		return fmt.Errorf("failed to create certificate file: %w", err)
	}
	defer certFile.Close()
	if err := pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes}); err != nil {
		return fmt.Errorf("failed to write certificate: %w", err)
	}

	// Save private key
	privKey, err := GenerateKey(false, false) // Example: Generate a new key
	if err != nil {
		return fmt.Errorf("failed to generate private key: %w", err)
	}
	privKeyBytes, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		return fmt.Errorf("failed to encode private key: %w", err)
	}
	keyFile, err := os.Create(keyPath)
	if err != nil {
		return fmt.Errorf("failed to create key file: %w", err)
	}
	defer keyFile.Close()
	if err := pem.Encode(keyFile, &pem.Block{Type: "PRIVATE KEY", Bytes: privKeyBytes}); err != nil {
		return fmt.Errorf("failed to write private key: %w", err)
	}

	return nil
}