package casp

import "github.com/persistenceOne/persistenceBridge/utilities/consterr"

const (
	ErrCASPSign             = consterr.Error("CASP SignOperation error")
	ErrCASPUncompressedKeys = consterr.Error("CASP error while getting UncompressedPublicKeys")
)
