package db

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"log"
	"time"
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
	if err == badger.ErrKeyNotFound {
		err = SetUnboundEpochTime(time.Now().Add(configuration.GetAppConfig().Kafka.EthUnbondCycleTime).Unix())
		if err != nil {
			return db, err
		}
	}
	// TODO add validator logic

	return db, nil
}

func OpenDB(dbPath string) (*badger.DB, error) {
	dbTemp, err := badger.Open(badger.DefaultOptions(dbPath))
	if err != nil {
		log.Fatalln(err)
	}
	db = dbTemp

	return db, nil
}
