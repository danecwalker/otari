package utils

import "io"

func WriteHeader(w io.Writer, header string) error {
	_, err := w.Write([]byte("[" + header + "]\n"))
	return err
}
