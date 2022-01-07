/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package db

import (
	"encoding/json"

	"github.com/dgraph-io/badger/v3"
	"github.com/ethereum/go-ethereum/common"
)

type IncomingEthereumTx struct {
	TxHash   common.Hash
	Sender   common.Address
	MsgBytes []byte
	MsgType  string
}

var _ KeyValue = &IncomingEthereumTx{}

func (t *IncomingEthereumTx) prefix() storeKeyPrefix {
	return incomingEthereumTxPrefix
}

func (t *IncomingEthereumTx) Key() []byte {
	return t.prefix().GenerateStoreKey(t.TxHash.Bytes())
}

func (t *IncomingEthereumTx) Value() ([]byte, error) {
	return json.Marshal(*t)
}

func (t *IncomingEthereumTx) Validate() error {
	if t.TxHash == EthEmptyHash() {
		return ErrEmptyTransaction
	}

	if len(t.MsgBytes) == 0 {
		return ErrEmptyTransactionMessage
	}

	if t.MsgType == "" {
		return ErrInvalidTransactionType
	}

	return nil
}

func GetIncomingEthereumTx(db *badger.DB, txHash common.Hash) (IncomingEthereumTx, error) {
	var ethInTx IncomingEthereumTx
	ethInTx.TxHash = txHash

	b, err := get(db, ethInTx.Key())
	if err != nil {
		return ethInTx, err
	}

	err = json.Unmarshal(b, &ethInTx)

	return ethInTx, err
}

func AddIncomingEthereumTx(db *badger.DB, t *IncomingEthereumTx) error {
	return set(db, t)
}

func CheckIncomingEthereumTxExists(db *badger.DB, txHash common.Hash) bool {
	ethInTx := IncomingEthereumTx{
		TxHash: txHash,
	}

	return keyExists(db, ethInTx.Key())
}
