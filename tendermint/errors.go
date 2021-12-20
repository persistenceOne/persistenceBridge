package tendermint

import (
	"github.com/persistenceOne/persistenceBridge/utilities/consterr"
)

const (
	ErrUnmarshalOutgoingTransaction = consterr.Error("failed to unmarshal OutgoingTendermintTransaction")
	ErrInvalidTxHash                = consterr.Error("invalid tx hash")
	ErrParseTransaction             = consterr.Error("unable to parse transaction into signing.Tx")
	ErrTransactionMessageGeneration = consterr.Error("failed to generate msgBytes")
	ErrAddToKafkaQueue              = consterr.Error("failed to add message(s) to kafka queue")
	ErrPartialSend                  = consterr.Error("added to IncomingTendermintTx but NOT to TendermintTxToKafka failed. Tx won't be added to kafka")
	ErrGetIncomingTendermintTx      = consterr.Error("failed to get IncomingTendermintTx by TendermintTxToKafka")
)
