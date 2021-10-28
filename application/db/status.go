package db

import (
	"encoding/json"
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

func getKafkaEthereumConsumeStatus() (Status, error) {
	return getStatus(kafkaEthConsume)
}

func setKafkaEthereumConsumeStatus(height int64) error {
	return setStatus(kafkaEthConsume, height)
}

func getKafkaTendermintConsumeStatus() (Status, error) {
	return getStatus(kafkaTMConsume)
}

func setKafkaTendermintConsumeStatus(height int64) error {
	return setStatus(kafkaTMConsume, height)
}
