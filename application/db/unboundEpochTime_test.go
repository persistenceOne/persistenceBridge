package db

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestSetUnboundEpochTime(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	require.Nil(t, err)

	var epochTime int64 = 4772132
	err = SetUnboundEpochTime(epochTime)
	require.Nil(t, err)

	db.Close()
}
func TestGetUnboundEpochTime(t *testing.T) {
	dirname, _ := os.UserHomeDir()
	db, err := OpenDB(filepath.Join(dirname, "/persistence/persistenceBridge/application") + "/db")
	require.Nil(t, err)

	var epochTime int64 = 4772132
	err = SetUnboundEpochTime(epochTime)
	require.Nil(t, err)

	var u UnboundEpochTime
	key := unboundEpochTimePrefix.GenerateStoreKey([]byte(unboundEpochTime))
	b, err := get(key)
	require.Nil(t, err)

	err = json.Unmarshal(b, &u)
	require.Nil(t, err)

	newUnboundEpochTime, err := GetUnboundEpochTime()
	require.Nil(t, err)

	require.Equal(t, reflect.TypeOf(epochTime), reflect.TypeOf(newUnboundEpochTime.Epoch))
	require.Equal(t, newUnboundEpochTime.Epoch, epochTime)

	db.Close()
}

func TestUnboundEpochTimeKey(t *testing.T) {
	var epochTime int64 = 4772132
	unboundEpochTime := UnboundEpochTime{
		Epoch: epochTime,
	}

	key := unboundEpochTime.Key()
	expectedKey := unboundEpochTimePrefix.GenerateStoreKey([]byte("UNBOUND_EPOCH_TIME"))

	require.Equal(t, reflect.TypeOf(key), reflect.TypeOf(expectedKey))
	require.Equal(t, expectedKey, key)
}

func TestUnboundEpochTimeValue(t *testing.T) {
	var epochTime int64 = 4772132
	unboundEpochTime := UnboundEpochTime{
		Epoch: epochTime,
	}
	expectedValue, err := json.Marshal(unboundEpochTime)
	require.Nil(t, err)

	value, err := unboundEpochTime.Value()
	require.Nil(t, err)

	require.Equal(t, reflect.TypeOf(value), reflect.TypeOf(expectedValue))
	require.Equal(t, expectedValue, value)
}

func TestUnboundEpochTimePrefix(t *testing.T) {
	var epochTime int64 = 4772132
	unboundEpochTime := UnboundEpochTime{
		Epoch: epochTime,
	}

	Prefix := unboundEpochTime.prefix()

	require.Equal(t, reflect.TypeOf(Prefix), reflect.TypeOf(unboundEpochTimePrefix))
	require.Equal(t, unboundEpochTimePrefix, Prefix)
}
