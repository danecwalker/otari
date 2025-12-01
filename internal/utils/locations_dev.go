//go:build !production

package utils

func OutputLocation() string {
	return "./stack"
}

func DataDirectory() string {
	return "./data_dev"
}
