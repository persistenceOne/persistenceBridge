/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package db

import (
	"encoding/binary"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestGenerateStoreKey(t *testing.T) {
	key := common.Hash{}.Bytes()
	storeKey := outgoingEthereumTxPrefix.GenerateStoreKey(key)
	Bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(Bytes, uint16(outgoingEthereumTxPrefix))

	Key := append(Bytes, key...)
	require.Equal(t, Key, storeKey)
}
