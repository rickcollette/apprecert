package config

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Config represents the application configuration.
type Config struct {
	CAROOT string
	CACert *x509.Certificate
}

// Load initializes and loads the configuration.
func Load() *Config {
	root := getCAROOT()
	return &Config{CAROOT: root}
}

// getCAROOT determines the default CA root directory.
func getCAROOT() string {
	if env := os.Getenv("CAROOT"); env != "" {
		return env
	}

	var dir string
	switch runtime.GOOS {
	case "windows":
		dir = filepath.Join(os.Getenv("LOCALAPPDATA"), "Apprecert")
	case "darwin":
		dir = filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "Apprecert")
	default:
		dir = filepath.Join(os.Getenv("HOME"), ".local", "share", "apprecert")
	}
	return dir
}

// LoadCA loads the CA certificate and key from the CAROOT.
func (cfg *Config) LoadCA() (*x509.Certificate, interface{}, error) {
	certPath := filepath.Join(cfg.CAROOT, "rootCA.pem")
	keyPath := filepath.Join(cfg.CAROOT, "rootCA-key.pem")

	// Load CA certificate
	certBytes, err := os.ReadFile(certPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read CA certificate: %w", err)
	}
	certBlock, _ := pem.Decode(certBytes)
	if certBlock == nil || certBlock.Type != "CERTIFICATE" {
		return nil, nil, fmt.Errorf("invalid CA certificate PEM")
	}
	caCert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse CA certificate: %w", err)
	}

	// Load CA key
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read CA key: %w", err)
	}
	keyBlock, _ := pem.Decode(keyBytes)
	if keyBlock == nil || keyBlock.Type != "PRIVATE KEY" {
		return nil, nil, fmt.Errorf("invalid CA key PEM")
	}
	caKey, err := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse CA key: %w", err)
	}

	return caCert, caKey, nil
}
