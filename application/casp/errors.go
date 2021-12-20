package casp

import "github.com/persistenceOne/persistenceBridge/utilities/consterr"

const (
	ErrNoPublicKeys              = consterr.Error("no public keys got from casp")
	ErrEmptyOperationID          = consterr.Error("empty operationID")
	ErrCantGetOperationApprovals = consterr.Error("unable to get approvals for operation")
)
