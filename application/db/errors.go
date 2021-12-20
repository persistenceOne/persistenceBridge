package db

import "github.com/persistenceOne/persistenceBridge/utilities/consterr"

const (
	ErrInvalidAmount  = consterr.Error("invalid amount")
	ErrNilAmount      = consterr.Error("amount is nil")
	ErrNegativeAmount = consterr.Error("amount less than 0")

	ErrEmptyTransaction        = consterr.Error("tx hash is empty")
	ErrEmptyTransactionMessage = consterr.Error("empty MsgBytes")
	ErrNoTransactionMessages   = consterr.Error("number of messages for ethHash is 0")

	ErrInvalidTransactionType = consterr.Error("invalid msg type")
	ErrEmptyDenom             = consterr.Error("empty denom")
	ErrEmptyFromAddress       = consterr.Error("from address empty")
	ErrInvalidFromAddress     = consterr.Error("invalid from address")
)
