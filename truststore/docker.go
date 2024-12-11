package truststore

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/appremon/apprecert/config"
)

func RemoveDockerTrust(cfg *config.Config) error {
	destPath := "/usr/local/share/ca-certificates/rootCA.crt"
	if err := os.Remove(destPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove Docker trust certificate: %w", err)
	}

	// Update CA certificates in Docker
	cmd := exec.Command("sudo", "update-ca-certificates")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to update Docker CA certificates: %w", err)
	}

	return nil
}
// UpdateDockerTrust copies the certificate to Docker images and updates trust.
func UpdateDockerTrust(cfg *config.Config) error {
	certPath := filepath.Join(cfg.CAROOT, "rootCA.pem")
	destPath := "/usr/local/share/ca-certificates/rootCA.crt"

	// Copy certificate to the appropriate directory
	cmd := exec.Command("sudo", "cp", certPath, destPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to copy certificate to Docker: %w", err)
	}

	// Update CA certificates in Docker
	cmd = exec.Command("sudo", "update-ca-certificates")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to update CA certificates: %w", err)
	}

	return nil
}
