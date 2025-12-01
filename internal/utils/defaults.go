package utils

import "os"

func DefaultStackPath() string {
	var defaultName = "stack.yaml"
	// check .yml if stack.yaml does not exist
	if fileExists(defaultName) {
		return defaultName
	}
	// return stack.yaml as default if neither exist
	return "stack.yml"
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func GetAbsolutePath(path string) (string, error) {
	if !isAbsolutePath(path) {
		absPath, err := os.Getwd()
		if err != nil {
			return "", err
		}
		return absPath + string(os.PathSeparator) + path, nil
	}
	return path, nil
}

func isAbsolutePath(path string) bool {
	return len(path) > 0 && path[0] == os.PathSeparator
}
