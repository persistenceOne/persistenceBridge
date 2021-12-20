package configuration

import (
	"github.com/persistenceOne/persistenceBridge/utilities/bridgeErr"
)

const (
	ErrIncorrectAccountPrefix = bridgeErr.Error("pStake wrap address prefix and config account prefix does not match")
	ErrInvalidBroadcastMode   = bridgeErr.Error("invalid broadcast mode")

	ErrRPCEndpointEmpty = bridgeErr.Error("rpc endpoint empty")

	ErrInvalidGasLimit    = bridgeErr.Error("invalid eth gas limit")
	ErrPStakeAddressEmpty = bridgeErr.Error("pStakeAddress empty")

	ErrEmptyAccountPrefix      = bridgeErr.Error("account prefix cannot be empty")
	ErrEmptyDenom              = bridgeErr.Error("denom cannot be empty")
	ErrNegativeWrapAmount      = bridgeErr.Error("minimum wrap amount cannot be less than 0")
	ErrEmptyChainID            = bridgeErr.Error("chain id cannot be empty")
	ErrInvalidTendermintNode   = bridgeErr.Error("invalid tendermint node")
	ErrTooLowReplicationFactor = bridgeErr.Error("replicationFactor has to be atleast 1")
	ErrTooFewParticipants      = bridgeErr.Error("num participants has to be atleast 1")
	ErrTooBigMinBatchSize      = bridgeErr.Error("min batch size cannot be greater than max batch size")

	ErrCaspVaultIDEmpty                  = bridgeErr.Error("casp vault id empty")
	ErrCaspAPITokenEmpty                 = bridgeErr.Error("casp api token empty")
	ErrCaspURLEmpty                      = bridgeErr.Error("casp url empty")
	ErrCaspPublicEmpty                   = bridgeErr.Error("casp ethereum public empty")
	ErrTooLowCaspMaxGetSignatureAttempts = bridgeErr.Error("casp MaxGetSignatureAttempts cannot be less than or equal to 0")

	ErrTelegramBotInvalidConfig = bridgeErr.Error("telegram bot configuration invalid")
)
