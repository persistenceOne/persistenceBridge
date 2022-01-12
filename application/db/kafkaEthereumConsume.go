package db

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/persistenceOne/persistenceBridge/application/constants"
)

type KafkaEthereumConsume struct {
	Index      uint64
	KafkaIndex int64
	MsgBytes   [][]byte
	TxHash     common.Hash
}

func (k KafkaEthereumConsume) prefix() storeKeyPrefix {
	return kafkaEthereumConsumePrefix
}

func (k KafkaEthereumConsume) Key() []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, k.Index)

	return k.prefix().GenerateStoreKey(bytes)
}

func (k KafkaEthereumConsume) Value() ([]byte, error) {
	return json.Marshal(k)
}

func (k KafkaEthereumConsume) Validate() error {
	if len(k.MsgBytes) == 0 {
		return fmt.Errorf("KafkaEthereumConsume: MsgBytes empty")
	}

	return nil
}

func AddKafkaEthereumConsume(database *badger.DB, kafkaIndex int64, msgBytes [][]byte) (uint64, error) {
	kafkaEthereumConsumeStatus, err := GetKafkaEthereumConsumeStatus(database)
	if err != nil {
		return 0, err
	}

	err = set(database, &KafkaEthereumConsume{
		Index:      uint64(kafkaEthereumConsumeStatus.LastCheckHeight),
		KafkaIndex: kafkaIndex,
		MsgBytes:   msgBytes,
		TxHash:     common.Hash{},
	})
	if err != nil {
		return uint64(kafkaEthereumConsumeStatus.LastCheckHeight), err
	}

	err = SetKafkaEthereumConsumeStatus(database, kafkaEthereumConsumeStatus.LastCheckHeight+1)
	if err != nil {
		return uint64(kafkaEthereumConsumeStatus.LastCheckHeight), err
	}

	return uint64(kafkaEthereumConsumeStatus.LastCheckHeight), err
}

func UpdateKafkaEthereumConsumeTxHash(database *badger.DB, index uint64, txHash common.Hash) error {
	k, err := GetKafkaEthereumConsume(database, index)
	if err != nil {
		return err
	}

	k.TxHash = txHash

	return set(database, &k)
}

func GetKafkaEthereumConsume(database *badger.DB, index uint64) (KafkaEthereumConsume, error) {
	result := KafkaEthereumConsume{
		Index: index,
	}

	b, err := get(database, result.Key())
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(b, &result)

	return result, err
}

func DeleteKafkaEthereumConsume(database *badger.DB, index uint64) error {
	kafkaEthereumConsume := KafkaEthereumConsume{
		Index: index,
	}

	return deleteKV(database, kafkaEthereumConsume.Key())
}

func GetEmptyTxHashesETH(database *badger.DB) ([]KafkaEthereumConsume, error) {
	var list []KafkaEthereumConsume

	err := iterateKeyValues(database, kafkaEthereumConsumePrefix.GenerateStoreKey([]byte{}), func(key []byte, value []byte) error {
		var k KafkaEthereumConsume

		err := json.Unmarshal(value, &k)
		if err != nil {
			return err
		}

		if k.TxHash == constants.EthereumEmptyTxHash() {
			list = append(list, k)
		}

		return nil
	})

	return list, err
}
