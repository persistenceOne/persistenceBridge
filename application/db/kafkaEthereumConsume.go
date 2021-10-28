package db

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
)

type KafkaEthereumConsume struct {
	Index      uint64
	KafkaIndex int64
	MsgBytes   []byte
	TxHash     common.Hash
}

var _ DBI = &KafkaEthereumConsume{}

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

func AddKafkaEthereumConsume(kafkaIndex int64, msgBytes []byte) error {
	kafkaEthereumConsumeStatus, err := getKafkaEthereumConsumeStatus()
	if err != nil {
		return err
	}

	err = set(&KafkaEthereumConsume{
		Index:      uint64(kafkaEthereumConsumeStatus.LastCheckHeight),
		KafkaIndex: kafkaIndex,
		MsgBytes:   msgBytes,
		TxHash:     common.Hash{},
	})
	if err != nil {
		return err
	}

	return setKafkaEthereumConsumeStatus(kafkaEthereumConsumeStatus.LastCheckHeight + 1)
}

func UpdateKafkaEthereumConsumeTxHash(index uint64, txHash common.Hash) error {
	k, err := GetKafkaEthereumConsume(index)
	if err != nil {
		return err
	}
	k.TxHash = txHash
	return set(&k)
}

func GetKafkaEthereumConsume(index uint64) (KafkaEthereumConsume, error) {
	var result KafkaEthereumConsume
	result.Index = index
	b, err := get(result.Key())
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(b, &result)
	return result, err
}
