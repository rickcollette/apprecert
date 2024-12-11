package truststore

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/appremon/apprecert/config"
)

func installLinux(cfg *config.Config) error {
	certPath := filepath.Join(cfg.CAROOT, "rootCA.pem")
	cmd := exec.Command("sudo", "cp", certPath, "/usr/local/share/ca-certificates/rootCA.crt")
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("sudo", "update-ca-certificates")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to update CA certificates: %s", out)
	}
	return nil
}

func uninstallLinux(_ *config.Config) error {
	cmd := exec.Command("sudo", "rm", "/usr/local/share/ca-certificates/rootCA.crt")
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("sudo", "update-ca-certificates")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to remove CA certificates: %s", out)
	}
	return nil
}
