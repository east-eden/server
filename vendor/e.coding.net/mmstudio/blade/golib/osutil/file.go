package osutil

import (
	"os"

	"golang.org/x/sys/unix"
)

// Exists function that determines if a given path exists.
func Exists(filePath string) (exists bool) {
	exists = true

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		exists = false
	}

	return exists
}

// IsWritable determines if a given directory or file can be written to.
func IsWritable(path string) (writable bool) {
	return unix.Access(path, unix.W_OK) == nil
}

// IsReadable determines if a given directory or file can be read from.
func IsReadable(path string) (readable bool) {
	return unix.Access(path, unix.R_OK) == nil
}
