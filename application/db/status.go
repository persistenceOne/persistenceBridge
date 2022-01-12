/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package db

import (
	"encoding/json"

	"github.com/dgraph-io/badger/v3"
)

const (
	cosmos          = "COSMOS"
	ethereum        = "ETHEREUM"
	kafkaEthConsume = "KAFKA_ETHEREUM_CONSUME"
	kafkaTMConsume  = "KAFKA_TENDERMINT_CONSUME"
)

type Status struct {
	Name            string
	LastCheckHeight int64
}

var _ KeyValue = &Status{}

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
	return nil
}

func getStatus(db *badger.DB, name string) (Status, error) {
	var status Status
	status.Name = name

	b, err := get(db, status.Key())
	if err != nil {
		return status, err
	}

	err = json.Unmarshal(b, &status)

	return status, err
}

func setStatus(db *badger.DB, name string, height int64) error {
	status := Status{
		Name:            name,
		LastCheckHeight: height,
	}

	return set(db, &status)
}

func GetCosmosStatus(db *badger.DB) (Status, error) {
	return getStatus(db, cosmos)
}

func SetCosmosStatus(db *badger.DB, height int64) error {
	return setStatus(db, cosmos, height)
}

func GetEthereumStatus(db *badger.DB) (Status, error) {
	return getStatus(db, ethereum)
}

func SetEthereumStatus(db *badger.DB, height int64) error {
	return setStatus(db, ethereum, height)
}

func GetKafkaEthereumConsumeStatus(db *badger.DB) (Status, error) {
	return getStatus(db, kafkaEthConsume)
}

func SetKafkaEthereumConsumeStatus(db *badger.DB, height int64) error {
	return setStatus(db, kafkaEthConsume, height)
}

func GetKafkaTendermintConsumeStatus(db *badger.DB) (Status, error) {
	return getStatus(db, kafkaTMConsume)
}

func SetKafkaTendermintConsumeStatus(db *badger.DB, height int64) error {
	return setStatus(db, kafkaTMConsume, height)
}
