package outgoingtx

import (
	"github.com/persistenceOne/persistenceBridge/utilities/consterr"
)

const (
	ErrNoWrapTokenMessages = consterr.Error("no wrap token messages to broadcast")
	ErrNoSignature         = consterr.Error("signature not found from casp for operation")
)
