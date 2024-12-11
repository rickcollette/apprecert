package truststore

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/appremon/apprecert/config"
)

// ConfigureGit sets Git to trust the generated certificate.
func ConfigureGit(cfg *config.Config) error {
	certPath := filepath.Join(cfg.CAROOT, "rootCA.pem")
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		return fmt.Errorf("certificate not found at %s", certPath)
	}

	cmd := exec.Command("git", "config", "--global", "http.sslCAInfo", certPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to configure Git trust: %w", err)
	}

	return nil
}

func UnconfigureGit(cfg *config.Config) error {
	cmd := exec.Command("git", "config", "--global", "--unset", "http.sslCAInfo")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to unconfigure Git trust: %w", err)
	}

	return nil
}