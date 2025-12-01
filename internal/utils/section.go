package utils

import "io"

func WriteSection(w io.Writer, section string, values [][2]string) error {
	err := WriteHeader(w, section)
	for _, v := range values {
		if err := WriteValue(w, v[0], v[1]); err != nil {
			return err
		}
	}
	return err
}
