package test

import (
	"os"
)

func TempDir() (dir string, closerFn func() error, err error) {
	closerFn = func() error { return nil }

	dir, err = os.MkdirTemp("", "*-badger")
	if err != nil {
		return
	}

	closerFn = func() error {
		return os.RemoveAll(dir)
	}

	return
}
