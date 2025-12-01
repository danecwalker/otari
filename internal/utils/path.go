package utils

import (
	"os"
	"path/filepath"
)

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

func StackNameFromPath(path string) string {
	filepathWithoutExt := filepath.Base(path)
	ext := filepath.Ext(filepathWithoutExt)
	return filepathWithoutExt[0 : len(filepathWithoutExt)-len(ext)]
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
