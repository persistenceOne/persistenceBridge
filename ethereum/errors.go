package ethereum

import "github.com/persistenceOne/persistenceBridge/utilities/bridgeErr"

const (
	ErrUnableGetMethodAndArgument = bridgeErr.Error("unable to get method and arguments")
	ErrFailedProcessArguments     = bridgeErr.Error("failed to process arguments of contract")
	ErrCannotAddTxToKafka         = bridgeErr.Error("added to IncomingEthereumTx but NOT to EthereumTxToKafka failed. Tx won't be added to kafka")
	ErrTxUnmarshal                = bridgeErr.Error("failed to unmarshal OutgoingEthereumTransaction")
)
