//go:build units

/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package db

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/persistenceOne/persistenceBridge/utilities/test"
	"github.com/stretchr/testify/require"
)

func TestSetStatus(t *testing.T) {
	database, closeFn, err := test.OpenDB(t, OpenDB)
	defer closeFn()

	require.Nil(t, err)

	const name = "tx1"

	lastCheckHeight := int64(4772132)

	err = setStatus(database, name, lastCheckHeight)
	require.Nil(t, err)
}

func TestGetStatus(t *testing.T) {
	database, closeFn, err := test.OpenDB(t, OpenDB)
	defer closeFn()

	require.Nil(t, err)

	name := "tx1"
	lastCheckHeight := int64(4772132)

	err = setStatus(database, name, lastCheckHeight)
	require.Nil(t, err)

	var expectedStatus Status
	expectedStatus.Name = name

	b, err := get(database, expectedStatus.Key())
	require.Nil(t, err)

	err = json.Unmarshal(b, &expectedStatus)
	require.Nil(t, err)

	status, err := getStatus(database, name)
	require.Nil(t, err)
	require.Equal(t, expectedStatus, status)
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
	database, closeFn, err := test.OpenDB(t, OpenDB)
	defer closeFn()

	require.Nil(t, err)

	err = SetCosmosStatus(database, 1)
	require.Nil(t, err)

	s, err := GetCosmosStatus(database)
	require.Nil(t, err)
	require.Equal(t, cosmos, s.Name)
	require.Equal(t, int64(1), s.LastCheckHeight)
}

func TestSetEthereumStatus(t *testing.T) {
	database, closeFn, err := test.OpenDB(t, OpenDB)
	defer closeFn()

	require.Nil(t, err)

	err = SetEthereumStatus(database, 1)
	require.Nil(t, err)

	s, err := GetEthereumStatus(database)
	require.Nil(t, err)
	require.Equal(t, ethereum, s.Name)
	require.Equal(t, int64(1), s.LastCheckHeight)
}
