//go:build units

/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package db

import (
	"encoding/json"
	"testing"

	"github.com/persistenceOne/persistenceBridge/utilities/test"
	"github.com/stretchr/testify/require"
)

func TestSetUnboundEpochTime(t *testing.T) {
	database, closeFn, err := test.OpenDB(t, OpenDB)
	defer closeFn()

	require.Nil(t, err)

	var epochTime int64 = 4772132
	err = SetUnboundEpochTime(database, epochTime)
	require.Nil(t, err)
}

func TestGetUnboundEpochTime(t *testing.T) {
	database, closeFn, err := test.OpenDB(t, OpenDB)
	defer closeFn()

	require.Nil(t, err)

	var epochTime int64 = 4772132
	err = SetUnboundEpochTime(database, epochTime)
	require.Nil(t, err)

	var u UnboundEpochTime
	key := unboundEpochTimePrefix.GenerateStoreKey([]byte(unboundEpochTime))

	b, err := get(database, key)
	require.Nil(t, err)

	err = json.Unmarshal(b, &u)
	require.Nil(t, err)

	newUnboundEpochTime, err := GetUnboundEpochTime(database)
	require.Nil(t, err)
	require.Equal(t, newUnboundEpochTime.Epoch, epochTime)
}

func TestUnboundEpochTimeKey(t *testing.T) {
	var epochTime int64 = 4772132
	unboundEpochTime := UnboundEpochTime{
		Epoch: epochTime,
	}

	key := unboundEpochTime.Key()
	expectedKey := unboundEpochTimePrefix.GenerateStoreKey([]byte("UNBOUND_EPOCH_TIME"))

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

	require.Equal(t, expectedValue, value)
}

func TestUnboundEpochTimePrefix(t *testing.T) {
	var epochTime int64 = 4772132
	unboundEpochTime := UnboundEpochTime{
		Epoch: epochTime,
	}

	Prefix := unboundEpochTime.prefix()

	require.Equal(t, unboundEpochTimePrefix, Prefix)
}
