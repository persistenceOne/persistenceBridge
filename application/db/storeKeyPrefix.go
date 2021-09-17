package db

import "encoding/binary"

type storeKeyPrefix int16

const (
	statusPrefix storeKeyPrefix = iota + 1
	validatorPrefix
	outgoingTendermintTxPrefix
	outgoingEthereumTxPrefix
	unboundEpochTimePrefix
	accountLimiterPrefix // Beta feature DO NOT REMOVE
	incomingTendermintTxPrefix
	incomingEthereumTxPrefix
	tendermintTxToKafkaPrefix
	ethereumTxToKafkaPrefix
)

func (storeKeyPrefix storeKeyPrefix) GenerateStoreKey(key []byte) []byte {
	bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(bytes, uint16(storeKeyPrefix))

	return append(bytes, key...)
}
