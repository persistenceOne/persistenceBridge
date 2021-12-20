/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package db

import (
	"errors"
	"log"
	"time"

	"github.com/dgraph-io/badger/v3"

	"github.com/persistenceOne/persistenceBridge/application/configuration"
)

var db *badger.DB

func InitializeDB(dbPath string, tendermintStart, ethereumStart int64) (*badger.DB, error) {
	dbTemp, err := badger.Open(badger.DefaultOptions(dbPath))
	if err != nil {
		log.Fatalln(err)
	}

	db = dbTemp

	if tendermintStart > 0 {
		err = SetCosmosStatus(tendermintStart - 1)
		if err != nil {
			return db, err
		}
	}

	if ethereumStart > 0 {
		err = SetEthereumStatus(ethereumStart - 1)
		if err != nil {
			return db, err
		}
	}

	_, err = GetUnboundEpochTime()
	if errors.Is(err, badger.ErrKeyNotFound) {
		err = SetUnboundEpochTime(time.Now().Add(configuration.GetAppConfig().Kafka.EthUnbondCycleTime).Unix())
		if err != nil {
			return db, err
		}
	}

	return db, nil
}

func OpenDB(dbPath string) (*badger.DB, error) {
	dbTemp, err := badger.Open(badger.DefaultOptions(dbPath))
	if err != nil {
		return nil, err
	}

	db = dbTemp

	return db, nil
}
