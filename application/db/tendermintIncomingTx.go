package db

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/dgraph-io/badger/v3"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
	tmBytes "github.com/tendermint/tendermint/libs/bytes"
)

type TendermintIncomingTx struct {
	TxHash          tmBytes.HexBytes
	MsgIndex        uint
	ProducedToKafka bool
	Msg             bankTypes.MsgSend
	Memo            string
}

var _ DBI = &TendermintIncomingTx{}

func (t *TendermintIncomingTx) prefix() storeKeyPrefix {
	return tendermintIncomingTxPrefix
}

func (t *TendermintIncomingTx) Key() []byte {
	producedToKafkaByte := []byte{byte(0)}
	if t.ProducedToKafka {
		producedToKafkaByte = []byte{byte(1)}
	}
	msgIndexBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(msgIndexBytes, uint16(t.MsgIndex))
	return t.prefix().GenerateStoreKey(append(producedToKafkaByte, append(t.TxHash, msgIndexBytes...)...))
}

func (t *TendermintIncomingTx) Value() ([]byte, error) {
	return json.Marshal(*t)
}

func (t *TendermintIncomingTx) Validate() error {
	if len(t.TxHash.Bytes()) == 0 {
		return fmt.Errorf("empty tx hash")
	}
	return nil
}

func GetTendermintIncomingTx(txHash tmBytes.HexBytes, msgIndex uint, producedToKafka bool) (TendermintIncomingTx, error) {
	var tmInTx TendermintIncomingTx
	tmInTx.TxHash = txHash
	tmInTx.MsgIndex = msgIndex
	tmInTx.ProducedToKafka = producedToKafka
	b, err := get(tmInTx.Key())
	if err != nil {
		return tmInTx, err
	}
	err = json.Unmarshal(b, &tmInTx)
	return tmInTx, err
}

func CheckTendermintIncomingTxProduced(txHash tmBytes.HexBytes, msgIndex uint) bool {
	var tmInTx TendermintIncomingTx
	tmInTx.TxHash = txHash
	tmInTx.MsgIndex = msgIndex
	tmInTx.ProducedToKafka = true
	result := true
	_, err := get(tmInTx.Key())
	if err != nil {
		if err == badger.ErrKeyNotFound {
			result = false
		} else {
			logging.Fatal(err)
		}
	}
	return result
}

func SetTendermintIncomingTx(a TendermintIncomingTx) error {
	return set(&a)
}

func GetProduceToKafkaTendermintTxs() ([]TendermintIncomingTx, error) {
	var result []TendermintIncomingTx
	err := iterateKeyValues(tendermintIncomingTxPrefix.GenerateStoreKey([]byte{byte(0)}), func(key []byte, value []byte) error {
		var t TendermintIncomingTx
		jsonErr := json.Unmarshal(value, &t)
		if jsonErr != nil {
			return jsonErr
		}
		result = append(result, t)
		return nil
	})
	return result, err
}
