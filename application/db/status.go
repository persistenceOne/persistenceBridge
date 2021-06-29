package db

import (
	"encoding/json"
	"github.com/dgraph-io/badger/v3"
)

const (
	COSMOS     = "COSMOS"
	ETHEREUM   = "ETHEREUM"
	VALIDATORS = "VALIDATORS"
)

type Status struct {
	Name            string
	LastCheckHeight int64
}

var _ DBI = &Status{}

func (status *Status) prefix() storeKeyPrefix {
	return statusPrefix
}

func (status *Status) Key() []byte {
	return status.prefix().GenerateStoreKey([]byte(status.Name))
}

func (status *Status) Value() ([]byte, error) {
	return json.Marshal(*status)
}

func (status *Status) Validate() error {
	// TODO
	return nil
}

func getStatus(name string) (Status, error) {
	var status Status
	status.Name = name
	b, err := get(status.Key())
	if err != nil {
		return Status{}, err
	}
	err = json.Unmarshal(b, &status)
	return status, err
}

func setStatus(name string, height int64) error {
	status := Status{
		Name:            name,
		LastCheckHeight: height,
	}
	return set(&status)
}

func GetCosmosStatus() (Status, error) {
	return getStatus(COSMOS)
}

func SetCosmosStatus(height int64) error {
	return setStatus(COSMOS, height)
}

func GetEthereumStatus() (Status, error) {
	return getStatus(ETHEREUM)
}

func SetEthereumStatus(height int64) error {
	return setStatus(ETHEREUM, height)
}

func GetValidators() ([]string, error) {
	var status []string
	err := db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(VALIDATORS))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			err = json.Unmarshal(val, &status)
			return err
		})
		return err
	})
	if err != nil {
		return status, err
	}
	return status, nil
}

func SetValidators(validators []string) error {
	b, err := json.Marshal(validators)
	if err != nil {
		return err
	}
	err = db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(VALIDATORS), b)
	})
	if err != nil {
		return err
	}
	return nil
}
