package casp

import "github.com/persistenceOne/persistenceBridge/utilities/bridgeErr"

const (
	ErrNoGroupsFound    = bridgeErr.Error("no groups found")
	ErrPendingApprovals = bridgeErr.Error("groups have pending approvals")
)
