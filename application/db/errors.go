package db

import "github.com/persistenceOne/persistenceBridge/utilities/bridgeErr"

const (
	ErrInvalidAmount  = bridgeErr.Error("invalid amount")
	ErrNilAmount      = bridgeErr.Error("amount is nil")
	ErrNegativeAmount = bridgeErr.Error("amount less than 0")

	ErrEmptyTransaction        = bridgeErr.Error("tx hash is empty")
	ErrEmptyTransactionMessage = bridgeErr.Error("empty MsgBytes")
	ErrNoTransactionMessages   = bridgeErr.Error("number of messages for ethHash is 0")

	ErrInvalidTransactionType = bridgeErr.Error("invalid msg type")
	ErrEmptyDenom             = bridgeErr.Error("empty denom")
	ErrEmptyFromAddress       = bridgeErr.Error("from address empty")
	ErrInvalidFromAddress     = bridgeErr.Error("invalid from address")
)
