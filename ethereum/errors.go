package ethereum

import "github.com/persistenceOne/persistenceBridge/utilities/consterr"

const (
	ErrUnableGetMethodAndArgument = consterr.Error("unable to get method and arguments")
	ErrFailedProcessArguments     = consterr.Error("failed to process arguments of contract")
	ErrCannotAddTxToKafka         = consterr.Error("added to IncomingEthereumTx but NOT to EthereumTxToKafka failed. Tx won't be added to kafka")
	ErrTxUnmarshal                = consterr.Error("failed to unmarshal OutgoingEthereumTransaction")
)
