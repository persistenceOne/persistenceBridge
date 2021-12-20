package commands

import "github.com/persistenceOne/persistenceBridge/utilities/bridgeErr"

const (
	ErrNoValidators = bridgeErr.Error("need to have at least one validator to redelegate to")
)
