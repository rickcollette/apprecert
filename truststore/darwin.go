package truststore

import (
	"os/exec"
	"path/filepath"

	"github.com/appremon/apprecert/config"
)

func installDarwin(cfg *config.Config) error {
	cmd := exec.Command("security", "add-trusted-cert", "-d", "-k", "/Library/Keychains/System.keychain", filepath.Join(cfg.CAROOT, "rootCA.pem"))
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func uninstallDarwin(cfg *config.Config) error {
	cmd := exec.Command("security", "remove-trusted-cert", "-d", filepath.Join(cfg.CAROOT, "rootCA.pem"))
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
