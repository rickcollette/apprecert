package cert

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"log"
	"time"

	"github.com/appremon/apprecert/config"
)

func Generate(cfg *config.Config, hosts []string) error {
    caCert, caKey, err := cfg.LoadCA()
    if err != nil {
        return fmt.Errorf("failed to load CA: %w", err)
    }

    // Assert that caKey implements crypto.Signer
    signer, ok := caKey.(crypto.Signer)
    if !ok {
        return fmt.Errorf("CA key does not implement crypto.Signer")
    }

	// Generate certificate template
	certTpl := &x509.Certificate{
		SerialNumber: randomSerialNumber(),
		Subject: pkix.Name{
			Organization: []string{"Apprecert Development Certificate"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(2, 3, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		BasicConstraintsValid: true,
	}

	// Generate certificate
	certBytes, err := x509.CreateCertificate(rand.Reader, certTpl, caCert, signer.Public(), signer)
	if err != nil {
		return err
	}

	// Save to disk
	err = saveCertificate(cfg, certBytes, hosts)
	if err != nil {
		return err
	}

	log.Printf("Certificate created for hosts: %v\n", hosts)
	return nil
}

