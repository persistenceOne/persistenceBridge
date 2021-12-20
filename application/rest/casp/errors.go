package casp

import "github.com/persistenceOne/persistenceBridge/utilities/bridgeErr"

const (
	ErrCASPSign             = bridgeErr.Error("CASP SignOperation error")
	ErrCASPUncompressedKeys = bridgeErr.Error("CASP error while getting UncompressedPublicKeys")
)
