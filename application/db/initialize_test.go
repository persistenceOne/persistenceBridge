//go:build units

/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package db

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/utilities/test"
)

func TestInitializeDB(t *testing.T) {
	configuration.SetConfig(test.GetCmdWithConfig())

	var (
		ethStart int64 = 4772131
		tmStart  int64 = 1
	)

	dir, dirCloser, err := test.TempDir()
	require.Nil(t, err)

	defer dirCloser()

	func() {
		database, err := InitializeDB(dir, tmStart, ethStart)
		defer database.Close()

		require.Nil(t, err)

		ethStatus, err := GetEthereumStatus(database)
		require.Nil(t, err)

		cosmosLastCheckHeight, err := GetCosmosStatus(database)
		require.Nil(t, err)

		ethHeight := ethStatus.LastCheckHeight + 1
		require.Equal(t, ethStart, ethHeight)

		cosmosHeight := cosmosLastCheckHeight.LastCheckHeight + 1
		require.Equal(t, tmStart, cosmosHeight)
	}()

	func() {
		database, err := OpenDB(dir)
		defer database.Close()

		require.Nil(t, err)

		err = deleteKV(database, unboundEpochTimePrefix.GenerateStoreKey([]byte(unboundEpochTime)))
		require.Nil(t, err)
	}()

	database, err := InitializeDB(dir, tmStart, ethStart)
	defer database.Close()

	require.Nil(t, err)
}

func TestOpenDB(t *testing.T) {
	_, closeFn, err := test.OpenDB(t, OpenDB)
	defer closeFn()

	require.Nil(t, err)

	_, closeFn1, err1 := test.OpenDB(t, OpenDB)
	defer closeFn1()

	require.Nil(t, err1)
}
