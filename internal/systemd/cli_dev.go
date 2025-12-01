//go:build !production

package systemd

func IsSystemdRunning() bool {
	return true
}

func IsUserLingeringEnabled() (bool, error) {
	return true, nil
}
