package db

import (
	"github.com/dgraph-io/badger/v3"
	"log"
)

var db *badger.DB

func InitializeDB(dbPath string, tendermintStart, ethereumStart int64, validators []string) (*badger.DB, error) {
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
