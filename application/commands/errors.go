package commands

import "github.com/persistenceOne/persistenceBridge/utilities/consterr"

const (
	ErrNoValidators = consterr.Error("need to have at least one validator to redelegate to")
)
