package utils

import "io"

func WriteEmptyLine(w io.Writer) error {
	_, err := w.Write([]byte("\n"))
	return err
}
