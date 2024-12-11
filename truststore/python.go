package truststore

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/appremon/apprecert/config"
)

// AppendToCertifi appends the certificate to the certifi bundle used by Python.
func AppendToCertifi(cfg *config.Config) error {
	certPath := filepath.Join(cfg.CAROOT, "rootCA.pem")
	certifiBundlePath := findCertifiBundle()
	if certifiBundlePath == "" {
		return fmt.Errorf("could not locate certifi bundle")
	}

	certBytes, err := os.ReadFile(certPath)
	if err != nil {
		return fmt.Errorf("failed to read certificate: %w", err)
	}

	bundleBytes, err := os.ReadFile(certifiBundlePath)
	if err != nil {
		return fmt.Errorf("failed to read certifi bundle: %w", err)
	}

	if bytes.Contains(bundleBytes, certBytes) {
		return nil // Already appended
	}

	err = os.WriteFile(certifiBundlePath, append(bundleBytes, certBytes...), 0644)
	if err != nil {
		return fmt.Errorf("failed to update certifi bundle: %w", err)
	}

	return nil
}

// findCertifiBundle locates the certifi bundle path.
func findCertifiBundle() string {
	pythonPath, err := exec.LookPath("python3")
	if err != nil {
		return ""
	}
	cmd := exec.Command(pythonPath, "-m", "certifi")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return string(bytes.TrimSpace(output))
}
