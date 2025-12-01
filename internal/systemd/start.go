package systemd

import "os/exec"

func StartUnit(unitName string) error {
	cmd := exec.Command("systemctl", "--user", "start", unitName)
	return cmd.Run()
}
