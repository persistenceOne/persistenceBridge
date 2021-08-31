package db

import "encoding/binary"

type storeKeyPrefix int16

const (
	statusPrefix storeKeyPrefix = iota + 1
	validatorPrefix
	tendermintBroadcastedTransactionPrefix
	ethereumBroadcastedWrapTokenTransactionPrefix
	unboundEpochTimePrefix
	accountLimiterPrefix // Beta feature DO NOT REMOVE
	tendermintIncomingTxPrefix
	ethereumIncomingTxPrefix
)

func (storeKeyPrefix storeKeyPrefix) GenerateStoreKey(key []byte) []byte {
	bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(bytes, uint16(storeKeyPrefix))

	return append(bytes, key...)
}
