package utils

import "os"

// PathExists checks if a file or directory exists.
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
