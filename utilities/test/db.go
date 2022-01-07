package test

import (
	"fmt"
	"testing"

	"github.com/dgraph-io/badger/v3"
	"github.com/stretchr/testify/require"
)

func OpenDB(t *testing.T, openDB func(dbPath string) (*badger.DB, error)) (*badger.DB, func(), error) {
	t.Helper()

	var (
		closers []func() error
		dir     string
	)

	dir, dirCloser, err := TempDir()
	if err != nil {
		return nil, func() {}, err
	}

	database, err := openDB(dir)
	if err != nil {
		return nil, func() {
			require.Nil(t, dirCloser())
		}, err
	}

	closers = append(closers, database.Close, dirCloser)

	closeFn := func() error {
		var innerErr error

		for _, fn := range closers {
			if fnErr := fn(); fnErr != nil {
				if innerErr == nil {
					innerErr = fn()
				} else {
					innerErr = fmt.Errorf("%w %v", fnErr, innerErr)
				}
			}
		}

		return innerErr
	}

	return database, func() {
		require.Nil(t, closeFn())
	}, nil
}
