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

func InitializeDB(dbPath string, tendermintStart, ethereumStart int64) (*badger.DB, error) {
	database, err := badger.Open(badger.DefaultOptions(dbPath))
	if err != nil {
		log.Fatalln(err)
	}

	if tendermintStart > 0 {
		err = SetCosmosStatus(database, tendermintStart-1)
		if err != nil {
			return nil, err
		}
	}

	if ethereumStart > 0 {
		err = SetEthereumStatus(database, ethereumStart-1)
		if err != nil {
			return nil, err
		}
	}

	_, err = GetUnboundEpochTime(database)
	if errors.Is(err, badger.ErrKeyNotFound) {
		err = SetUnboundEpochTime(database, time.Now().Add(configuration.GetAppConfig().Kafka.EthUnbondCycleTime).Unix())
		if err != nil {
			return nil, err
		}
	}

	return database, nil
}

func OpenDB(dbPath string) (*badger.DB, error) {
	return badger.Open(badger.DefaultOptions(dbPath))
}
