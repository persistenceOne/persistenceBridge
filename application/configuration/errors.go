package configuration

import (
	"github.com/persistenceOne/persistenceBridge/utilities/consterr"
)

const (
	ErrIncorrectAccountPrefix = consterr.Error("pStake wrap address prefix and config account prefix does not match")
	ErrInvalidBroadcastMode   = consterr.Error("invalid broadcast mode")

	ErrRPCEndpointEmpty = consterr.Error("rpc endpoint empty")

	ErrInvalidGasLimit    = consterr.Error("invalid eth gas limit")
	ErrPStakeAddressEmpty = consterr.Error("pStakeAddress empty")

	ErrEmptyAccountPrefix      = consterr.Error("account prefix cannot be empty")
	ErrEmptyDenom              = consterr.Error("denom cannot be empty")
	ErrNegativeWrapAmount      = consterr.Error("minimum wrap amount cannot be less than 0")
	ErrEmptyChainID            = consterr.Error("chain id cannot be empty")
	ErrInvalidTendermintNode   = consterr.Error("invalid tendermint node")
	ErrTooLowReplicationFactor = consterr.Error("replicationFactor has to be atleast 1")
	ErrTooFewParticipants      = consterr.Error("num participants has to be atleast 1")
	ErrTooBigMinBatchSize      = consterr.Error("min batch size cannot be greater than max batch size")

	ErrCaspVaultIDEmpty                  = consterr.Error("casp vault id empty")
	ErrCaspAPITokenEmpty                 = consterr.Error("casp api token empty")
	ErrCaspURLEmpty                      = consterr.Error("casp url empty")
	ErrCaspPublicEmpty                   = consterr.Error("casp ethereum public empty")
	ErrTooLowCaspMaxGetSignatureAttempts = consterr.Error("casp MaxGetSignatureAttempts cannot be less than or equal to 0")

	ErrTelegramBotInvalidConfig = consterr.Error("telegram bot configuration invalid")
)
