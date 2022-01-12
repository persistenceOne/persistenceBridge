package ethtypes

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/persistenceOne/persistenceBridge/utilities/consterr"
	"github.com/spf13/cobra"
)

const (
	ErrInvalidAddress = consterr.Error("invalid address specified")
)

func FromStringFlag(cmd *cobra.Command, flagName string) (common.Address, error) {
	address, err := cmd.Flags().GetString(flagName)
	if err != nil {
		return common.Address{}, err
	}

	return FromString(address)
}

func FromString(addr string) (common.Address, error) {
	if !common.IsHexAddress(addr) {
		return common.Address{}, ErrInvalidAddress
	}

	return common.HexToAddress(addr), nil
}
