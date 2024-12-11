package truststore

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/appremon/apprecert/config"
)

// NSSProfile represents an NSS database profile.
type NSSProfile struct {
	Path        string
	CertutilCmd string
}

// FindNSSProfiles locates available NSS profiles on the system.
func FindNSSProfiles() ([]NSSProfile, error) {
	var profiles []string
	switch runtime.GOOS {
	case "linux":
		profiles = []string{
			filepath.Join(os.Getenv("HOME"), ".mozilla/firefox/*"),
			filepath.Join(os.Getenv("HOME"), ".pki/nssdb"),
		}
	case "darwin":
		profiles = []string{
			filepath.Join(os.Getenv("HOME"), "Library/Application Support/Firefox/Profiles/*"),
		}
	case "windows":
		profiles = []string{
			filepath.Join(os.Getenv("USERPROFILE"), "AppData\\Roaming\\Mozilla\\Firefox\\Profiles\\*"),
		}
	default:
		return nil, fmt.Errorf("unsupported platform for NSS profiles")
	}

	var foundProfiles []NSSProfile
	for _, pattern := range profiles {
		matches, _ := filepath.Glob(pattern)
		for _, path := range matches {
			if _, err := os.Stat(filepath.Join(path, "cert9.db")); err == nil {
				foundProfiles = append(foundProfiles, NSSProfile{Path: path, CertutilCmd: "certutil"})
			}
		}
	}
	return foundProfiles, nil
}

// Install adds the CA certificate to the NSS trust store.
func (n *NSSProfile) Install(cfg *config.Config) error {
	certPath := filepath.Join(cfg.CAROOT, "rootCA.pem")
	cmd := exec.Command(n.CertutilCmd, "-A", "-d", "sql:"+n.Path,
		"-t", "C,,", "-n", "apprecert-rootCA", "-i", certPath,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install certificate in NSS trust store (%s): %s", n.Path, output)
	}
	return nil
}

// Uninstall removes the CA certificate from the NSS trust store.
func (n *NSSProfile) Uninstall() error {
	cmd := exec.Command(n.CertutilCmd, "-D", "-d", "sql:"+n.Path, "-n", "apprecert-rootCA")
	output, err := cmd.CombinedOutput()
	if err != nil && !strings.Contains(string(output), "could not be found") {
		return fmt.Errorf("failed to remove certificate from NSS trust store (%s): %s", n.Path, output)
	}
	return nil
}
