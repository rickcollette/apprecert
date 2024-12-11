package truststore

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/appremon/apprecert/config"
)

// ConfigureNodeJS sets NODE_EXTRA_CA_CERTS to include the generated certificate.
func ConfigureNodeJS(cfg *config.Config) error {
	certPath := filepath.Join(cfg.CAROOT, "rootCA.pem")
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		return fmt.Errorf("certificate not found at %s", certPath)
	}

	os.Setenv("NODE_EXTRA_CA_CERTS", certPath)
	return nil
}

func UnconfigureNodeJS(cfg *config.Config) error {
	// Clear the NODE_EXTRA_CA_CERTS environment variable
	if err := os.Unsetenv("NODE_EXTRA_CA_CERTS"); err != nil {
		return fmt.Errorf("failed to unset NODE_EXTRA_CA_CERTS: %w", err)
	}

	return nil
}