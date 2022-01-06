/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package db

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/persistenceOne/persistenceBridge/application/constants"
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
	if t.TxHash.String() == constants.EthereumEmptyTxHash {
		return fmt.Errorf("tx hash is empty")
	}
	if len(t.MsgBytes) == 0 {
		return fmt.Errorf("empty MsgBytes")
	}
	if t.MsgType == "" {
		return fmt.Errorf("invalid msg type")
	}
	if t.Sender.String() == constants.EthereumZeroAddress {
		return fmt.Errorf("invalid sender address")
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
