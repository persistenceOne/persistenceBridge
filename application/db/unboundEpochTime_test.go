/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package db

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistenceBridge/application/constants"
)

func TestSetUnboundEpochTime(t *testing.T) {
	database, err := OpenDB(constants.TestDbDir)
	require.Nil(t, err)

	var epochTime int64 = 4772132
	err = SetUnboundEpochTime(epochTime)
	require.Nil(t, err)

	database.Close()
}
func TestGetUnboundEpochTime(t *testing.T) {
	database, err := OpenDB(constants.TestDbDir)
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

	database.Close()
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
