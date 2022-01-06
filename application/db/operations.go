/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package db

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

func get(key []byte) ([]byte, error) {
	var dbi []byte
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			dbi = val
			return nil
		})
		return err
	})
	return dbi, err
}

func keyExists(key []byte) bool {
	err := db.View(func(txn *badger.Txn) error {
		_, err := txn.Get(key)
		return err
	})
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return false
		} else {
			logging.Fatal(err)
		}
	}
	return true
}

func set(dbi DBI) error {
	err := dbi.Validate()
	if err != nil {
		return err
	}
	b, err := dbi.Value()
	if err != nil {
		return err
	}
	err = db.Update(func(txn *badger.Txn) error {
		return txn.Set(dbi.Key(), b)
	})
	if err != nil {
		return err
	}
	return nil
}

func iterateKeyValues(prefix []byte, operation func(key []byte, value []byte) error) error {
	return db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(value []byte) error {
				return operation(item.Key(), value)
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// iterateKeys doesn't fetch item values until called upon by item.Value()
func iterateKeys(prefix []byte, operation func(key []byte, item *badger.Item) error) error {
	return db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := operation(item.Key(), item)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func deleteKV(key []byte) error {
	return db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}
