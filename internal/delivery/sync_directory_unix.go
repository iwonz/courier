//go:build !windows

package delivery

import (
	"errors"
	"os"
)

func syncStateDirectory(root *os.Root) error {
	directory, err := root.Open(".")
	if err != nil {
		return err
	}
	return errors.Join(directory.Sync(), directory.Close())
}
