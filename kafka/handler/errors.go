package handler

import "github.com/persistenceOne/persistenceBridge/utilities/consterr"

const (
	ErrKafkaNilMessage       = consterr.Error("kafka returned nil message")
	ErrKafkaErrorMessage     = consterr.Error("kafka returned error message")
	ErrTooFewDelegatedTokens = consterr.Error("unbondings requested are greater than delegated tokens")
)
