package rpc

import "github.com/persistenceOne/persistenceBridge/utilities/bridgeErr"

const (
	ErrHTTPWriter = bridgeErr.Error("can't write an error into HTTP response")
)
