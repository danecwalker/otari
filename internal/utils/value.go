package utils

import "io"

func WriteValue(w io.Writer, key, value string) error {
	_, err := w.Write([]byte(key + "=" + value + "\n"))
	return err
}
