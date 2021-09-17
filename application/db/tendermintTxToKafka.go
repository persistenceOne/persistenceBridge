package db

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dgraph-io/badger/v3"
	tmBytes "github.com/tendermint/tendermint/libs/bytes"
)

type TendermintTxToKafka struct {
	TxHash   tmBytes.HexBytes
	MsgIndex uint
	Denom    string
}

var _ DBI = &TendermintTxToKafka{}

func (t *TendermintTxToKafka) prefix() storeKeyPrefix {
	return tendermintTxToKafkaPrefix
}

func (t *TendermintTxToKafka) Key() []byte {
	msgIndexBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(msgIndexBytes, uint16(t.MsgIndex))
	denomBytes := []byte(t.Denom)
	return t.prefix().GenerateStoreKey(append(t.TxHash, append(msgIndexBytes, denomBytes...)...))
}

func (t *TendermintTxToKafka) Value() ([]byte, error) {
	return []byte{}, nil
}

func (t *TendermintTxToKafka) Validate() error {
	if len(t.TxHash.Bytes()) == 0 {
		return fmt.Errorf("empty tx hash")
	}
	if t.Denom == "" {
		return fmt.Errorf("empty denom")
	}
	if err := sdk.ValidateDenom(t.Denom); err != nil {
		return err
	}
	return nil
}

func AddTendermintTxToKafka(t TendermintTxToKafka) error {
	return set(&t)
}

func GetAllTendermintTxToKafka() ([]TendermintTxToKafka, error) {
	var tmTxToKafkaList []TendermintTxToKafka
	err := iterateKeys(tendermintTxToKafkaPrefix.GenerateStoreKey([]byte{}), func(key []byte, _ *badger.Item) error {
		var t TendermintTxToKafka
		err := json.Unmarshal(key, &t)
		if err != nil {
			return err
		}
		tmTxToKafkaList = append(tmTxToKafkaList, t)
		return nil
	})
	return tmTxToKafkaList, err
}

func DeleteTendermintTxToKafka(txHash tmBytes.HexBytes, msgIndex uint, denom string) error {
	tmInTx := TendermintTxToKafka{
		TxHash:   txHash,
		MsgIndex: msgIndex,
		Denom:    denom,
	}
	return deleteKV(tmInTx.Key())
}
