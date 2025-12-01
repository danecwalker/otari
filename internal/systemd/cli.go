package systemd

import (
	"fmt"
	"os"
	"os/user"
)

func IsSystemdRunning() bool {
	// A common way to check if systemd is running is to check for the presence of the
	// /run/systemd/system directory, which is created when systemd is active.
	if _, err := os.Stat("/run/systemd/system"); os.IsNotExist(err) {
		return false
	}
	return true
}

func IsUserLingeringEnabled() (bool, error) {
	// Check for the existence of the lingering file for the current user
	// username $USER
	user, err := user.Current()
	if err != nil {
		return false, err
	}

	lingeringFilePath := fmt.Sprintf("/var/lib/systemd/linger/%s", user.Username)

	if _, err := os.Stat(lingeringFilePath); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}
