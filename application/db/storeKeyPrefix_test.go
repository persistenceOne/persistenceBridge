package db

import (
	"encoding/binary"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenerateStoreKey(t *testing.T) {
	key := common.Hash{}.Bytes()
	storeKey := outgoingEthereumTxPrefix.GenerateStoreKey(key)
	Bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(Bytes, uint16(outgoingEthereumTxPrefix))

	Key := append(Bytes, key...)
	require.Equal(t, Key, storeKey)
}
