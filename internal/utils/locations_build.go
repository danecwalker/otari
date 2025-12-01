//go:build production

package utils

import (
	"os"
	"path/filepath"
)

func OutputLocation() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "stack"
	}
	return filepath.Join(homeDir, ".config", "containers", "systemd")
}
