/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package db

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
)

type IncomingEthereumTx struct {
	TxHash   common.Hash
	Sender   common.Address
	MsgBytes []byte
	MsgType  string
}

var _ DBI = &IncomingEthereumTx{}

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

func GetIncomingEthereumTx(txHash common.Hash) (IncomingEthereumTx, error) {
	var ethInTx IncomingEthereumTx
	ethInTx.TxHash = txHash
	b, err := get(ethInTx.Key())
	if err != nil {
		return ethInTx, err
	}
	err = json.Unmarshal(b, &ethInTx)
	return ethInTx, err
}

func AddIncomingEthereumTx(t IncomingEthereumTx) error {
	return set(&t)
}

func CheckIncomingEthereumTxExists(txHash common.Hash) bool {
	ethInTx := IncomingEthereumTx{
		TxHash: txHash,
	}
	return keyExists(ethInTx.Key())
}
