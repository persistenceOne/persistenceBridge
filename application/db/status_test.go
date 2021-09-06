package db

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestSetStatus(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")

	Name := "tx1"
	var LastCheckHeight int64 = 4772132
	err = setStatus(Name, LastCheckHeight)
	require.Nil(t, err)

	db.Close()
}

func TestGetStatus(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")

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

	db.Close()

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

	Prefix := status.prefix()

	require.Equal(t, reflect.TypeOf(statusPrefix), reflect.TypeOf(Prefix))
	require.Equal(t, statusPrefix, Prefix)
}
