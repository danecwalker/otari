package systemd

import "os/exec"

func ReloadDaemon() error {
	cmd := exec.Command("systemctl", "--user", "daemon-reload")
	return cmd.Run()
}
