//go:build !windows
// +build !windows

package truststore

import "github.com/appremon/apprecert/config"

func installWindows(_ *config.Config) error {
	// No-op on non-Windows systems
	return nil
}

func uninstallWindows(_ *config.Config) error {
	// No-op on non-Windows systems
	return nil
}