/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package db

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	test "github.com/persistenceOne/persistenceBridge/utilities/testing"
)

func TestInitializeDB(t *testing.T) {
	configuration.InitConfig()
	configuration.SetConfig(test.GetCmdWithConfig())

	var (
		ethStart int64 = 4772131
		tmStart  int64 = 1
	)

	database, err := InitializeDB(constants.TestHomeDir, tmStart, ethStart)
	require.Nil(t, err)

	ethStatus, err := GetEthereumStatus()
	require.Nil(t, err)

	cosmosLastCheckHeight, err := GetCosmosStatus()
	require.Nil(t, err)

	ethHeight := ethStatus.LastCheckHeight + 1
	require.Equal(t, ethStart, ethHeight)

	cosmosHeight := cosmosLastCheckHeight.LastCheckHeight + 1
	require.Equal(t, tmStart, cosmosHeight)

	database.Close()

	database, err = OpenDB(constants.TestHomeDir)

	err = deleteKV(unboundEpochTimePrefix.GenerateStoreKey([]byte(unboundEpochTime)))
	require.Nil(t, err)

	database.Close()

	database, err = InitializeDB(constants.TestHomeDir, tmStart, ethStart)
	require.Nil(t, err)

	database.Close()
}

func TestOpenDB(t *testing.T) {
	db, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	newDb, err := OpenDB(constants.TestDbDir)
	require.Nil(t, newDb)
	require.Equal(t, "Cannot acquire directory lock on \""+constants.TestDbDir+"\".  Another process is using this Badger database. error: resource temporarily unavailable", err.Error())

	db.Close()
}
