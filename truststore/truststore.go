package truststore

import (
	"runtime"

	"github.com/appremon/apprecert/config"
)

func Install(cfg *config.Config) error {
	// Platform-specific installations
	if runtime.GOOS == "darwin" {
		if err := installDarwin(cfg); err != nil {
			return err
		}
	} else if runtime.GOOS == "linux" {
		if err := installLinux(cfg); err != nil {
			return err
		}
	} else if runtime.GOOS == "windows" {
		if err := installWindows(cfg); err != nil {
			return err
		}
	}

	// Kubernetes trust store
	if err := InstallKubernetes(cfg); err != nil {
		return err
	}

	// Docker trust store
	if err := UpdateDockerTrust(cfg); err != nil {
		return err
	}

	// Git trust store
	if err := ConfigureGit(cfg); err != nil {
		return err
	}

	// Node.js trust store
	if err := ConfigureNodeJS(cfg); err != nil {
		return err
	}

	// Java trust store
	javaStore, err := NewJavaTrustStore()
	if err == nil {
		if err := javaStore.Install(cfg); err != nil {
			return err
		}
	}

	// NSS profiles
	nssProfiles, err := FindNSSProfiles()
	if err == nil {
		for _, profile := range nssProfiles {
			if err := profile.Install(cfg); err != nil {
				return err
			}
		}
	}

	return nil
}

func Uninstall(cfg *config.Config) error {
	// Platform-specific uninstallations
	if runtime.GOOS == "darwin" {
		if err := uninstallDarwin(cfg); err != nil {
			return err
		}
	} else if runtime.GOOS == "linux" {
		if err := uninstallLinux(cfg); err != nil {
			return err
		}
	} else if runtime.GOOS == "windows" {
		if err := uninstallWindows(cfg); err != nil {
			return err
		}
	}

	// Kubernetes trust store
	if err := UninstallKubernetes(cfg); err != nil {
		return err
	}

	// Docker trust store
	if err := RemoveDockerTrust(cfg); err != nil { // You may need to create a `RemoveDockerTrust` function.
		return err
	}

	// Git trust store
	if err := UnconfigureGit(cfg); err != nil { // You may need to create a `UnconfigureGit` function.
		return err
	}

	// Node.js trust store
	if err := UnconfigureNodeJS(cfg); err != nil { // You may need to create a `UnconfigureNodeJS` function.
		return err
	}

	// Java trust store
	javaStore, err := NewJavaTrustStore()
	if err == nil {
		if err := javaStore.Uninstall(); err != nil {
			return err
		}
	}

	// NSS profiles
	nssProfiles, err := FindNSSProfiles()
	if err == nil {
		for _, profile := range nssProfiles {
			if err := profile.Uninstall(); err != nil {
				return err
			}
		}
	}

	return nil
}
