package tendermint

import (
	"github.com/persistenceOne/persistenceBridge/utilities/bridgeErr"
)

const (
	ErrUnmarshalOutgoingTransaction = bridgeErr.Error("failed to unmarshal OutgoingTendermintTransaction")
	ErrInvalidTxHash                = bridgeErr.Error("invalid tx hash")
	ErrParseTransaction             = bridgeErr.Error("unable to parse transaction into signing.Tx")
	ErrTransactionMessageGeneration = bridgeErr.Error("failed to generate msgBytes")
	ErrAddToKafkaQueue              = bridgeErr.Error("failed to add message(s) to kafka queue")
	ErrPartialSend                  = bridgeErr.Error("added to IncomingTendermintTx but NOT to TendermintTxToKafka failed. Tx won't be added to kafka")
	ErrGetIncomingTendermintTx      = bridgeErr.Error("failed to get IncomingTendermintTx by TendermintTxToKafka")
)
