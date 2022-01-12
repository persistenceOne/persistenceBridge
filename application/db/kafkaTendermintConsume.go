package db

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/badger/v3"
	tmBytes "github.com/tendermint/tendermint/libs/bytes"
)

type KafkaTendermintConsume struct {
	Index      uint64
	KafkaIndex int64
	MsgBytes   [][]byte
	TxHash     tmBytes.HexBytes
}

func (k KafkaTendermintConsume) prefix() storeKeyPrefix {
	return kafkaTendermintConsumePrefix
}

func (k KafkaTendermintConsume) Key() []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, k.Index)

	return k.prefix().GenerateStoreKey(bytes)
}

func (k KafkaTendermintConsume) Value() ([]byte, error) {
	return json.Marshal(k)
}

func (k KafkaTendermintConsume) Validate() error {
	if len(k.MsgBytes) == 0 {
		return fmt.Errorf("KafkaTendermintConsume: MsgBytes empty")
	}

	return nil
}

func AddKafkaTendermintConsume(database *badger.DB, kafkaIndex int64, msgBytes [][]byte) (uint64, error) {
	kafkaTMConsumeStatus, err := GetKafkaTendermintConsumeStatus(database)
	if err != nil {
		return 0, err
	}

	err = set(database, &KafkaTendermintConsume{
		Index:      uint64(kafkaTMConsumeStatus.LastCheckHeight),
		KafkaIndex: kafkaIndex,
		MsgBytes:   msgBytes,
		TxHash:     tmBytes.HexBytes{},
	})

	if err != nil {
		return 0, err
	}

	err = SetKafkaTendermintConsumeStatus(database, kafkaTMConsumeStatus.LastCheckHeight+1)
	if err != nil {
		return uint64(kafkaTMConsumeStatus.LastCheckHeight), err
	}

	return uint64(kafkaTMConsumeStatus.LastCheckHeight), err
}

func UpdateKafkaTendermintConsumeTxHash(database *badger.DB, index uint64, txHash tmBytes.HexBytes) error {
	k, err := GetKafkaTendermintConsume(database, index)
	if err != nil {
		return err
	}

	k.TxHash = txHash

	return set(database, &k)
}

func GetKafkaTendermintConsume(database *badger.DB, index uint64) (KafkaTendermintConsume, error) {
	result := KafkaTendermintConsume{
		Index: index,
	}

	b, err := get(database, result.Key())
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(b, &result)

	return result, err
}

func DeleteKafkaTendermintConsume(database *badger.DB, index uint64) error {
	kafkaTendermintConsume := KafkaTendermintConsume{
		Index: index,
	}

	return deleteKV(database, kafkaTendermintConsume.Key())
}

func GetEmptyTxHashesTM(database *badger.DB) ([]KafkaTendermintConsume, error) {
	var list []KafkaTendermintConsume

	err := iterateKeyValues(database, kafkaTendermintConsumePrefix.GenerateStoreKey([]byte{}), func(key []byte, value []byte) error {
		var k KafkaTendermintConsume

		innerErr := json.Unmarshal(value, &k)
		if innerErr != nil {
			return innerErr
		}

		if len(k.TxHash.Bytes()) == 0 {
			list = append(list, k)
		}

		return nil
	})
	return list, err
}
