package casp

import "github.com/persistenceOne/persistenceBridge/utilities/bridgeErr"

const (
	ErrNoPublicKeys              = bridgeErr.Error("no public keys got from casp")
	ErrEmptyOperationID          = bridgeErr.Error("empty operationID")
	ErrCantGetOperationApprovals = bridgeErr.Error("unable to get approvals for operation")
)
