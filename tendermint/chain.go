package tendermint

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceBridge/application/casp"
	"time"

	"github.com/cosmos/relayer/helpers"
	"github.com/cosmos/relayer/relayer"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
	tendermintService "github.com/tendermint/tendermint/libs/service"
)

func InitializeAndStartChain(timeout, homePath string) (*relayer.Chain, error) {
	chain := &relayer.Chain{}
	chain.Key = "unusedKey"
	chain.ChainID = configuration.GetAppConfig().Tendermint.ChainID
	chain.RPCAddr = configuration.GetAppConfig().Tendermint.Node
	chain.AccountPrefix = configuration.GetAppConfig().Tendermint.AccountPrefix
	chain.GasAdjustment = 1.5
	chain.GasPrices = "0.025" + configuration.GetAppConfig().Tendermint.PStakeDenom
	chain.TrustingPeriod = "21h"

	to, err := time.ParseDuration(timeout)
	if err != nil {
		return chain, err
	}

	err = chain.Init(homePath, to, nil, true)
	if err != nil {
		return chain, err
	}

	if chain.KeyExists(chain.Key) {
		logging.Info("deleting old key", chain.Key)
		err = chain.Keybase.Delete(chain.Key)
		if err != nil {
			return chain, err
		}
	}

	_, err = helpers.KeyAddOrRestore(chain, chain.Key, configuration.GetAppConfig().Tendermint.CoinType)
	if err != nil {
		return chain, err
	}

	if err = chain.Start(); err != nil {
		if err != tendermintService.ErrAlreadyStarted {
			chain.Error(err)
			return chain, err
		}
	}
	return chain, nil
}

func SetBech32PrefixesAndPStakeWrapAddress() (sdkTypes.AccAddress, error) {
	if configuration.GetAppConfig().Tendermint.AccountPrefix == "" {
		panic("account prefix is empty")
	}
	bech32PrefixAccAddr := configuration.GetAppConfig().Tendermint.AccountPrefix
	bech32PrefixAccPub := configuration.GetAppConfig().Tendermint.AccountPrefix + sdkTypes.PrefixPublic
	bech32PrefixValAddr := configuration.GetAppConfig().Tendermint.AccountPrefix + sdkTypes.PrefixValidator + sdkTypes.PrefixOperator
	bech32PrefixValPub := configuration.GetAppConfig().Tendermint.AccountPrefix + sdkTypes.PrefixValidator + sdkTypes.PrefixOperator + sdkTypes.PrefixPublic
	bech32PrefixConsAddr := configuration.GetAppConfig().Tendermint.AccountPrefix + sdkTypes.PrefixValidator + sdkTypes.PrefixConsensus
	bech32PrefixConsPub := configuration.GetAppConfig().Tendermint.AccountPrefix + sdkTypes.PrefixValidator + sdkTypes.PrefixConsensus + sdkTypes.PrefixPublic

	bech32Configuration := sdkTypes.GetConfig()
	bech32Configuration.SetBech32PrefixForAccount(bech32PrefixAccAddr, bech32PrefixAccPub)
	bech32Configuration.SetBech32PrefixForValidator(bech32PrefixValAddr, bech32PrefixValPub)
	bech32Configuration.SetBech32PrefixForConsensusNode(bech32PrefixConsAddr, bech32PrefixConsPub)
	// Do not seal the config.

	tmAddress, err := casp.GetTendermintAddress()
	if err != nil {
		return sdkTypes.AccAddress{}, err
	}

	configuration.SetPStakeAddress(tmAddress)

	return tmAddress, nil
}
