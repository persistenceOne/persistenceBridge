package config

import (
	"crypto/ecdsa"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

type Configuration struct {
	PStakeDenom          string
	PStakeAddress        sdkTypes.AccAddress
	EthAccountPrivateKey *ecdsa.PrivateKey
	EthGasLimit          uint64
}

var appConfiguration *Configuration

func SetAppConfiguration(denom string, address sdkTypes.AccAddress, ethAccountPrivateKey *ecdsa.PrivateKey, ethGasLimit uint64) {
	appConfiguration = &Configuration{
		PStakeDenom:          denom,
		PStakeAddress:        address,
		EthAccountPrivateKey: ethAccountPrivateKey,
		EthGasLimit:          ethGasLimit,
	}
}

func GetAppConfiguration() *Configuration {
	return appConfiguration
}
