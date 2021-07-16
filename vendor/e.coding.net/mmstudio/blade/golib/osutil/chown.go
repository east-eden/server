// +build !linux

package osutil

import (
	"os"
)

func Chown(_ string, _ os.FileInfo) error {
	return nil
}
