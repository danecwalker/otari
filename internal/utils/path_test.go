package utils

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAbsolutePath(t *testing.T) {
	tests := []struct {
		input    string
		expected string // expected absolute path
	}{
		{
			"relative/path/to/file",
			filepath.Join("/mnt", "d", "projects", "otari", "internal", "utils", "relative", "path", "to", "file"),
		},
		{"/absolute/path/to/file", "/absolute/path/to/file"},
		{"", filepath.Join("/mnt", "d", "projects", "otari", "internal", "utils")},
		{".", filepath.Join("/mnt", "d", "projects", "otari", "internal", "utils")},
		{"..", filepath.Join("/mnt", "d", "projects", "otari", "internal")},
		{"./", filepath.Join("/mnt", "d", "projects", "otari", "internal", "utils")},
	}

	for _, tt := range tests {
		absPath, err := GetAbsolutePath(tt.input)
		assert.NoError(t, err)
		assert.NotEmpty(t, absPath)
		assert.Equal(t, tt.expected, absPath, "Absolute path does not match expected value")
	}
}
