//go:build production

package systemd

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/danecwalker/otari/internal/utils"
)

func ReloadDaemon() error {
	cmd := exec.Command("systemctl", "--user", "daemon-reload")
	return cmd.Run()
}

func StartUnit(unitName string) error {
	cmd := exec.Command("systemctl", "--user", "start", unitName)
	return cmd.Run()
}

func StopUnit(unitName string) error {
	cmd := exec.Command("systemctl", "--user", "stop", unitName)
	return cmd.Run()
}

func RestartUnit(unitName string) error {
	cmd := exec.Command("systemctl", "--user", "restart", unitName)
	return cmd.Run()
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
