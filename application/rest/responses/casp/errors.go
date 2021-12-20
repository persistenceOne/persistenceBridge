package casp

import "github.com/persistenceOne/persistenceBridge/utilities/consterr"

const (
	ErrNoGroupsFound    = consterr.Error("no groups found")
	ErrPendingApprovals = consterr.Error("groups have pending approvals")
)
