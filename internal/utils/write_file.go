package utils

import "os"

func WriteFileAtomic(path string, data []byte, perm os.FileMode) error {
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, perm); err != nil {
		return err
	}
	return os.Rename(tmpPath, path)
}

func WriteToFile(dir, filename string, data []byte) error {
	fullPath := dir + string(os.PathSeparator) + filename
	return WriteFileAtomic(fullPath, data, 0644)
}
