package truststore

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/appremon/apprecert/config"
)

// InstallKubernetes updates the Kubernetes cluster to trust the custom CA.
func InstallKubernetes(cfg *config.Config) error {
	certPath := filepath.Join(cfg.CAROOT, "rootCA.pem")
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		return fmt.Errorf("certificate not found at %s", certPath)
	}

	// Check if the ConfigMap already exists
	checkCmd := exec.Command("kubectl", "get", "configmap", "custom-ca-bundle", "-n", "kube-system")
	if err := checkCmd.Run(); err == nil {
		return fmt.Errorf("ConfigMap 'custom-ca-bundle' already exists in namespace 'kube-system'")
	}

	// Create a ConfigMap for the CA bundle
	cmd := exec.Command("kubectl", "create", "configmap", "custom-ca-bundle",
		"--from-file=ca.crt="+certPath, "-n", "kube-system")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create ConfigMap: %s, %w", output, err)
	}

	// Patch the cluster to include the ConfigMap
	patchCmd := exec.Command("kubectl", "patch", "cm", "kubeadm-config", "-n", "kube-system",
		"--type=json", "-p",
		"[{\"op\": \"add\", \"path\": \"/data/ClusterConfiguration/certificatesDir\", \"value\": \"/etc/kubernetes/pki/custom-ca-bundle\"}]")
	if output, err := patchCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to patch cluster configuration: %s, %w", output, err)
	}

	return nil
}

// UninstallKubernetes removes the custom CA from the Kubernetes cluster.
func UninstallKubernetes(cfg *config.Config) error {
	// Remove the ConfigMap
	cmd := exec.Command("kubectl", "delete", "configmap", "custom-ca-bundle", "-n", "kube-system")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to delete ConfigMap: %s, %w", output, err)
	}

	return nil
}
