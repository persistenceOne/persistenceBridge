package handler

import "github.com/persistenceOne/persistenceBridge/utilities/bridgeErr"

const (
	ErrKafkaNilMessage       = bridgeErr.Error("kafka returned nil message")
	ErrKafkaErrorMessage     = bridgeErr.Error("kafka returned error message")
	ErrTooFewDelegatedTokens = bridgeErr.Error("unbondings requested are greater than delegated tokens")
)
