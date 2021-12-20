package outgoingTx

import (
	"github.com/persistenceOne/persistenceBridge/utilities/bridgeErr"
)

const (
	ErrNoWrapTokenMessages = bridgeErr.Error("no wrap token messages to broadcast")
	ErrNoSignature         = bridgeErr.Error("signature not found from casp for operation")
)
