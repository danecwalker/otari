//go:build !production

package systemd

import (
	"os"
	"path/filepath"

	"github.com/danecwalker/otari/internal/utils"
)

func ReloadDaemon() error {
	return nil
}

func StartUnit(unitName string) error {
	return nil
}

func StopUnit(unitName string) error {
	return nil
}

func RestartUnit(unitName string) error {
	return nil
}

func DeleteUnitFile(unitName string) error {
	outputDir := utils.OutputLocation()
	unitFilePath := filepath.Join(outputDir, unitName)
	// Remove the unit file
	err := os.Remove(unitFilePath)
	if err != nil {
		// If the file does not exist, consider it successful
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return nil
}
