/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package db

import (
	"encoding/json"

	"github.com/dgraph-io/badger/v3"
)

const unboundEpochTime = "UNBOUND_EPOCH_TIME"

type UnboundEpochTime struct {
	Epoch int64
}

var _ KeyValue = &UnboundEpochTime{}

func (u *UnboundEpochTime) prefix() storeKeyPrefix {
	return unboundEpochTimePrefix
}

func (u *UnboundEpochTime) Key() []byte {
	return u.prefix().GenerateStoreKey([]byte(unboundEpochTime))
}

func (u *UnboundEpochTime) Value() ([]byte, error) {
	return json.Marshal(*u)
}

func (u *UnboundEpochTime) Validate() error {
	return nil
}

func GetUnboundEpochTime(database *badger.DB) (UnboundEpochTime, error) {
	var u UnboundEpochTime

	key := unboundEpochTimePrefix.GenerateStoreKey([]byte(unboundEpochTime))

	b, err := get(database, key)
	if err != nil {
		return u, err
	}

	err = json.Unmarshal(b, &u)

	return u, err
}

func SetUnboundEpochTime(database *badger.DB, epochTime int64) error {
	u := UnboundEpochTime{
		Epoch: epochTime,
	}

	return set(database, &u)
}
