package db

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dgraph-io/badger/v3"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
	tmBytes "github.com/tendermint/tendermint/libs/bytes"
)

type TendermintIncomingTx struct {
	ProducedToKafka bool
	TxHash          tmBytes.HexBytes
	MsgIndex        uint
	Denom           string
	FromAddress     string
	Amount          sdk.Int
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
	denomBytes := []byte(t.Denom)
	return t.prefix().GenerateStoreKey(append(producedToKafkaByte, append(t.TxHash, append(msgIndexBytes, denomBytes...)...)...))
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

func GetTendermintIncomingTx(txHash tmBytes.HexBytes, msgIndex uint, producedToKafka bool, denom string) (TendermintIncomingTx, error) {
	var tmInTx TendermintIncomingTx
	tmInTx.ProducedToKafka = producedToKafka
	tmInTx.TxHash = txHash
	tmInTx.MsgIndex = msgIndex
	tmInTx.Denom = denom
	b, err := get(tmInTx.Key())
	if err != nil {
		return tmInTx, err
	}
	err = json.Unmarshal(b, &tmInTx)
	return tmInTx, err
}

func CheckTendermintIncomingTxProduced(txHash tmBytes.HexBytes, msgIndex uint, denom string) bool {
	var tmInTx TendermintIncomingTx
	tmInTx.ProducedToKafka = true
	tmInTx.TxHash = txHash
	tmInTx.MsgIndex = msgIndex
	tmInTx.Denom = denom
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

func AddToPendingTendermintIncomingTx(t TendermintIncomingTx) error {
	return set(&t)
}

func SetTendermintIncomingTxProduced(t TendermintIncomingTx) error {
	t.ProducedToKafka = false
	err := deleteKV(t.Key())
	if err != nil {
		return err
	}
	t.ProducedToKafka = true
	return set(&t)
}

func GetProduceToKafkaTendermintIncomingTxs() ([]TendermintIncomingTx, error) {
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
