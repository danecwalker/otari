package utils

import "os"

func DefaultStackPath() string {
	var defaultName = "otari.yaml"
	// check .yml if otari.yaml does not exist
	if fileExists(defaultName) {
		return defaultName
	}
	// return otari.yaml as default if neither exist
	return "otari.yml"
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
