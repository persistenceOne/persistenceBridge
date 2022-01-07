/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package db

import (
	"encoding/json"
	"fmt"
	"math/big"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	tmBytes "github.com/tendermint/tendermint/libs/bytes"
)

type WrapTokenMsg struct {
	FromAddress      sdkTypes.AccAddress
	TendermintTxHash tmBytes.HexBytes
	Address          common.Address
	StakingAmount    *big.Int
	WrapAmount       *big.Int
}

type OutgoingEthereumTransaction struct {
	TxHash   common.Hash
	Messages []WrapTokenMsg
}

func (w WrapTokenMsg) Validate() error {
	if w.FromAddress.String() == "" {
		return fmt.Errorf("from address empty")
	}
	if w.Address.String() == constants.EthereumZeroAddress {
		return fmt.Errorf("invalid eth address")
	}
	if len(w.TendermintTxHash.Bytes()) != 32 {
		return fmt.Errorf("invalid tm tx hash")
	}
	if w.StakingAmount == nil {
		return fmt.Errorf("staking amount is nil")
	}
	if w.WrapAmount == nil {
		return fmt.Errorf("wrapping amount is nil")
	}
	if w.WrapAmount.String() == sdkTypes.ZeroInt().BigInt().String() && w.StakingAmount.String() == sdkTypes.ZeroInt().BigInt().String() {
		return fmt.Errorf("both amounts zero")
	}
	return nil
}

func NewWrapTokenMsg(fromAddress sdkTypes.AccAddress, tmTxHash tmBytes.HexBytes, stakingAmount *big.Int, ethAddress common.Address, wrapAmount *big.Int) WrapTokenMsg {
	return WrapTokenMsg{
		FromAddress:      fromAddress,
		TendermintTxHash: tmTxHash,
		Address:          ethAddress,
		StakingAmount:    stakingAmount,
		WrapAmount:       wrapAmount,
	}
}

func NewOutgoingETHTransaction(txHash common.Hash, msgs []WrapTokenMsg) OutgoingEthereumTransaction {
	return OutgoingEthereumTransaction{TxHash: txHash, Messages: msgs}
}

var _ DBI = &OutgoingEthereumTransaction{}

func (ethTx *OutgoingEthereumTransaction) prefix() storeKeyPrefix {
	return outgoingEthereumTxPrefix
}

func (ethTx *OutgoingEthereumTransaction) Key() []byte {
	return ethTx.prefix().GenerateStoreKey(ethTx.TxHash.Bytes())
}

func (ethTx *OutgoingEthereumTransaction) Value() ([]byte, error) {
	return json.Marshal(*ethTx)
}

func (ethTx *OutgoingEthereumTransaction) Validate() error {
	if ethTx.TxHash.String() == constants.EthereumEmptyTxHash {
		return fmt.Errorf("tx hash is empty")
	}
	if len(ethTx.Messages) == 0 {
		return fmt.Errorf("number of messages for ethHash %s is 0", ethTx.TxHash)
	}

	for _, msg := range ethTx.Messages {
		err := msg.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteOutgoingEthereumTx(txHash common.Hash) error {
	return deleteKV(outgoingEthereumTxPrefix.GenerateStoreKey(txHash.Bytes()))
}

func SetOutgoingEthereumTx(ethTransaction OutgoingEthereumTransaction) error {
	return set(&ethTransaction)
}

func IterateOutgoingEthTx(operation func(key []byte, value []byte) error) error {
	return iterateKeyValues(outgoingEthereumTxPrefix.GenerateStoreKey([]byte{}), operation)
}
