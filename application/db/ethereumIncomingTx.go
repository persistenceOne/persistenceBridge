package db

import (
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

type EthereumIncomingTx struct {
	TxHash          common.Hash
	ProducedToKafka bool
	Sender          common.Address
	MsgBytes        []byte
	MsgType         string
}

var _ DBI = &EthereumIncomingTx{}

func (t *EthereumIncomingTx) prefix() storeKeyPrefix {
	return ethereumIncomingTxPrefix
}

func (t *EthereumIncomingTx) Key() []byte {
	producedToKafkaByte := []byte{byte(0)}
	if t.ProducedToKafka {
		producedToKafkaByte = []byte{byte(1)}
	}
	return t.prefix().GenerateStoreKey(append(producedToKafkaByte, t.TxHash.Bytes()...))
}

func (t *EthereumIncomingTx) Value() ([]byte, error) {
	return json.Marshal(*t)
}

func (t *EthereumIncomingTx) Validate() error {
	if t.TxHash.String() == "0x0000000000000000000000000000000000000000000000000000000000000000" {
		return fmt.Errorf("tx hash is empty")
	}
	if len(t.MsgBytes) == 0 {
		return fmt.Errorf("empty MsgBytes")
	}
	if t.MsgType == "" {
		return fmt.Errorf("invalid msg type")
	}
	return nil
}

func GetEthereumIncomingTx(txHash common.Hash, producedToKafka bool) (EthereumIncomingTx, error) {
	var ethInTx EthereumIncomingTx
	ethInTx.TxHash = txHash
	ethInTx.ProducedToKafka = producedToKafka
	b, err := get(ethInTx.Key())
	if err != nil {
		return ethInTx, err
	}
	err = json.Unmarshal(b, &ethInTx)
	return ethInTx, err
}

func AddToPendingEthereumIncomingTx(t EthereumIncomingTx) error {
	t.ProducedToKafka = false
	return set(&t)
}

func SetEthereumIncomingTxProduced(t EthereumIncomingTx) error {
	t.ProducedToKafka = false
	err := deleteKV(t.Key())
	if err != nil {
		return err
	}
	t.ProducedToKafka = true
	return set(&t)
}

func GetProduceToKafkaEthereumIncomingTxs() ([]EthereumIncomingTx, error) {
	var result []EthereumIncomingTx
	err := iterateKeyValues(ethereumIncomingTxPrefix.GenerateStoreKey([]byte{byte(0)}), func(key []byte, value []byte) error {
		var e EthereumIncomingTx
		jsonErr := json.Unmarshal(value, &e)
		if jsonErr != nil {
			return jsonErr
		}
		result = append(result, e)
		return nil
	})
	return result, err
}

func CheckEthereumIncomingTxProduced(txHash common.Hash) bool {
	var ethInTx EthereumIncomingTx
	ethInTx.TxHash = txHash
	ethInTx.ProducedToKafka = true
	result := true
	_, err := get(ethInTx.Key())
	if err != nil {
		if err == badger.ErrKeyNotFound {
			result = false
		} else {
			logging.Fatal(err)
		}
	}
	return result
}
