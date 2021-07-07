package db

import (
	"encoding/json"
)

const (
	cosmos   = "COSMOS"
	ethereum = "ETHEREUM"
)

type Status struct {
	Name            string
	LastCheckHeight int64 //TODO change it to Index
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
		return status, err
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
	return getStatus(cosmos)
}

func SetCosmosStatus(height int64) error {
	return setStatus(cosmos, height)
}

func GetEthereumStatus() (Status, error) {
	return getStatus(ethereum)
}

func SetEthereumStatus(height int64) error {
	return setStatus(ethereum, height)
}
