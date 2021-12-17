/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package db

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmBytes "github.com/tendermint/tendermint/libs/bytes"
)

type IncomingTendermintTx struct {
	TxHash      tmBytes.HexBytes
	MsgIndex    uint
	Denom       string
	FromAddress string
	Amount      sdk.Int
	Memo        string
}

var _ DBI = &IncomingTendermintTx{}

func (t *IncomingTendermintTx) prefix() storeKeyPrefix {
	return incomingTendermintTxPrefix
}

func (t *IncomingTendermintTx) Key() []byte {
	msgIndexBytes := make([]byte, 2)

	binary.LittleEndian.PutUint16(msgIndexBytes, uint16(t.MsgIndex))

	denomBytes := []byte(t.Denom)

	return t.prefix().GenerateStoreKey(append(t.TxHash, append(msgIndexBytes, denomBytes...)...))
}

func (t *IncomingTendermintTx) Value() ([]byte, error) {
	return json.Marshal(*t)
}

func (t *IncomingTendermintTx) Validate() error {
	if len(t.TxHash.Bytes()) == 0 {
		return fmt.Errorf("empty tx hash")
	}

	if t.Denom == "" {
		return fmt.Errorf("empty denom")
	}

	if err := sdk.ValidateDenom(t.Denom); err != nil {
		return err
	}

	if t.FromAddress == "" {
		return fmt.Errorf("from address empty")
	}

	_, err := sdk.AccAddressFromBech32(t.FromAddress)
	if err != nil {
		return fmt.Errorf("invalid from address")
	}

	if t.Amount.IsNil() {
		return fmt.Errorf("amount is nil")
	}

	if t.Amount.LT(sdk.ZeroInt()) {
		return fmt.Errorf("amount less than 0")
	}

	return nil
}

func GetIncomingTendermintTx(txHash tmBytes.HexBytes, msgIndex uint, denom string) (IncomingTendermintTx, error) {
	var tmInTx IncomingTendermintTx

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

func AddIncomingTendermintTx(t IncomingTendermintTx) error {
	return set(&t)
}

func CheckIncomingTendermintTxExists(txHash tmBytes.HexBytes, msgIndex uint, denom string) bool {
	tmInTx := IncomingTendermintTx{
		TxHash:   txHash,
		MsgIndex: msgIndex,
		Denom:    denom,
	}

	return keyExists(tmInTx.Key())
}
