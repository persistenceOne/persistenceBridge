package db

import (
	"encoding/json"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestSetStatus(t *testing.T) {
	database, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	Name := "tx1"
	var LastCheckHeight int64 = 4772132
	err = setStatus(Name, LastCheckHeight)
	require.Nil(t, err)

	database.Close()
}

func TestGetStatus(t *testing.T) {
	database, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	Name := "tx1"
	var LastCheckHeight int64 = 4772132
	err = setStatus(Name, LastCheckHeight)
	require.Nil(t, err)

	var expectedStatus Status
	expectedStatus.Name = Name
	b, err := get(expectedStatus.Key())
	require.Nil(t, err)

	err = json.Unmarshal(b, &expectedStatus)

	status, err := getStatus(Name)
	require.Nil(t, err)
	require.Equal(t, expectedStatus, status)

	database.Close()

}

func TestStatusKey(t *testing.T) {
	status := Status{}
	expectedKey := status.prefix().GenerateStoreKey([]byte(status.Name))
	key := status.Key()

	require.Equal(t, expectedKey, key)
}

func TestStatusValue(t *testing.T) {
	status := Status{}

	expectedValue, err := json.Marshal(status)
	require.Nil(t, err)

	value, err := status.Value()
	require.Nil(t, err)

	require.Equal(t, expectedValue, value)

}

func TestStatusPrefix(t *testing.T) {
	status := Status{}
	require.Equal(t, reflect.TypeOf(statusPrefix), reflect.TypeOf(status.prefix()))
	require.Equal(t, statusPrefix, status.prefix())
}

func TestSetCosmosStatus(t *testing.T) {
	database, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	err = SetCosmosStatus(1)
	require.Nil(t, err)

	s, err := GetCosmosStatus()
	require.Nil(t, err)
	require.Equal(t, cosmos, s.Name)
	require.Equal(t, int64(1), s.LastCheckHeight)

	database.Close()
}

func TestSetEthereumStatus(t *testing.T) {
	database, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	err = SetEthereumStatus(1)
	require.Nil(t, err)

	s, err := GetEthereumStatus()
	require.Nil(t, err)
	require.Equal(t, ethereum, s.Name)
	require.Equal(t, int64(1), s.LastCheckHeight)

	database.Close()
}
