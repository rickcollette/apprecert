package cert

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"

	"github.com/appremon/apprecert/config"
	"software.sslmate.com/src/go-pkcs12"
)

// GenerateMultipleFormats saves certificates in .pem, .crt, and .p12 formats.
func GenerateMultipleFormats(cfg *config.Config, certBytes []byte, privKey crypto.PrivateKey) error {
	// Save PEM format
	pemPath := filepath.Join(cfg.CAROOT, "rootCA.pem")
	if err := savePEM(certBytes, pemPath); err != nil {
		return err
	}

	// Save CRT format
	crtPath := filepath.Join(cfg.CAROOT, "rootCA.crt")
	if err := savePEM(certBytes, crtPath); err != nil {
		return err
	}

	// Save P12 format
	p12Path := filepath.Join(cfg.CAROOT, "rootCA.p12")
	if err := saveP12(certBytes, privKey, p12Path); err != nil {
		return err
	}

	return nil
}

// savePEM saves the certificate in PEM format.
func savePEM(certBytes []byte, path string) error {
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	return os.WriteFile(path, certPEM, 0644)
}

// saveP12 saves the certificate and private key in P12 format.
func saveP12(certBytes []byte, privKey crypto.PrivateKey, path string) error {
	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return err
	}
	p12Data, err := pkcs12.Encode(rand.Reader, privKey, cert, nil, "changeit")
	if err != nil {
		return fmt.Errorf("failed to generate PKCS#12: %w", err)
	}
	return os.WriteFile(path, p12Data, 0644)
}
