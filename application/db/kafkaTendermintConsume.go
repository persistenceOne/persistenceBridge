package db

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	tmBytes "github.com/tendermint/tendermint/libs/bytes"
)

type KafkaTendermintConsume struct {
	Index      uint64
	KafkaIndex int64
	MsgBytes   [][]byte
	TxHash     tmBytes.HexBytes
}

var _ DBI = &KafkaTendermintConsume{}

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

func AddKafkaTendermintConsume(kafkaIndex int64, msgBytes [][]byte) (uint64, error) {
	kafkaTMConsumeStatus, err := GetKafkaTendermintConsumeStatus()
	if err != nil {
		return 0, err
	}

	err = set(&KafkaTendermintConsume{
		Index:      uint64(kafkaTMConsumeStatus.LastCheckHeight),
		KafkaIndex: kafkaIndex,
		MsgBytes:   msgBytes,
		TxHash:     tmBytes.HexBytes{},
	})
	err = SetKafkaTendermintConsumeStatus(kafkaTMConsumeStatus.LastCheckHeight + 1)
	if err != nil {
		return uint64(kafkaTMConsumeStatus.LastCheckHeight), err
	}
	return uint64(kafkaTMConsumeStatus.LastCheckHeight), err
}

func UpdateKafkaTendermintConsumeTxHash(index uint64, txHash tmBytes.HexBytes) error {
	k, err := GetKafkaTendermintConsume(index)
	if err != nil {
		return err
	}
	k.TxHash = txHash
	return set(&k)
}

func GetKafkaTendermintConsume(index uint64) (KafkaTendermintConsume, error) {
	var result KafkaTendermintConsume
	result.Index = index
	b, err := get(result.Key())
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(b, &result)
	return result, err
}

func DeleteKafkaTendermintConsume(index uint64) error {
	kafkaTendermintConsume := KafkaTendermintConsume{
		Index: index,
	}
	return deleteKV(kafkaTendermintConsume.Key())
}

func GetEmptyTxHashesTM() ([]KafkaTendermintConsume, error) {
	var list []KafkaTendermintConsume
	err := iterateKeyValues(kafkaTendermintConsumePrefix.GenerateStoreKey([]byte{}), func(key []byte, value []byte) error {
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
