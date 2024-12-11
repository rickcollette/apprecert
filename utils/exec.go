package utils

import (
	"os"
	"os/exec"
)

// CommandWithSudo returns an exec.Cmd with sudo prepended if not running as root.
func CommandWithSudo(cmd ...string) *exec.Cmd {
	if os.Geteuid() == 0 {
		return exec.Command(cmd[0], cmd[1:]...)
	}
	return exec.Command("sudo", cmd...)
}
