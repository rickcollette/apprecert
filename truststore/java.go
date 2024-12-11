package truststore

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/appremon/apprecert/config"
)

// Default values for Java trust store
const (
	defaultStorePass = "changeit" // Default password for cacerts
	cacertsFileName  = "cacerts"
)

// JavaTrustStore contains Java trust store configuration.
type JavaTrustStore struct {
	keytoolPath string
	cacertsPath string
	storePass   string
}

// NewJavaTrustStore initializes a new JavaTrustStore.
func NewJavaTrustStore() (*JavaTrustStore, error) {
	javaHome := os.Getenv("JAVA_HOME")
	if javaHome == "" {
		return nil, fmt.Errorf("JAVA_HOME is not set")
	}

	keytoolPath := filepath.Join(javaHome, "bin", "keytool")
	cacertsPath := filepath.Join(javaHome, "lib", "security", cacertsFileName)

	if _, err := os.Stat(keytoolPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("keytool not found in JAVA_HOME")
	}
	if _, err := os.Stat(cacertsPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("cacerts file not found in JAVA_HOME")
	}

	return &JavaTrustStore{
		keytoolPath: keytoolPath,
		cacertsPath: cacertsPath,
		storePass:   defaultStorePass,
	}, nil
}

// Install adds the CA certificate to the Java trust store.
func (j *JavaTrustStore) Install(cfg *config.Config) error {
	certPath := filepath.Join(cfg.CAROOT, "rootCA.pem")
	cmd := exec.Command(j.keytoolPath, "-importcert", "-noprompt",
		"-keystore", j.cacertsPath,
		"-storepass", j.storePass,
		"-file", certPath,
		"-alias", "apprecert-rootCA",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install certificate in Java trust store: %s", output)
	}
	return nil
}

// Uninstall removes the CA certificate from the Java trust store.
func (j *JavaTrustStore) Uninstall() error {
	cmd := exec.Command(j.keytoolPath, "-delete",
		"-keystore", j.cacertsPath,
		"-storepass", j.storePass,
		"-alias", "apprecert-rootCA",
	)
	output, err := cmd.CombinedOutput()
	if err != nil && !bytes.Contains(output, []byte("does not exist")) {
		return fmt.Errorf("failed to remove certificate from Java trust store: %s", output)
	}
	return nil
}
