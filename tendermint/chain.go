/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package tendermint

import (
	"context"
	"errors"
	"time"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/persistenceBridge/application/casp"

	"github.com/cosmos/relayer/helpers"
	"github.com/cosmos/relayer/relayer"
	tendermintService "github.com/tendermint/tendermint/libs/service"

	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

func getChain() *relayer.Chain {
	return &relayer.Chain{
		Key:            "unusedKey",
		ChainID:        configuration.GetAppConfig().Tendermint.ChainID,
		RPCAddr:        configuration.GetAppConfig().Tendermint.Node,
		AccountPrefix:  configuration.GetAppConfig().Tendermint.AccountPrefix,
		GasAdjustment:  1.5,
		GasPrices:      "0.025" + configuration.GetAppConfig().Tendermint.PStakeDenom,
		TrustingPeriod: "21h",
	}
}

func ChainInit(timeout, homePath string) (*relayer.Chain, error) {
	chain := getChain()

	to, err := time.ParseDuration(timeout)
	if err != nil {
		return nil, err
	}

	err = chain.Init(homePath, to, nil, true)
	if err != nil {
		return nil, err
	}

	return chain, nil
}

func InitializeAndStartChain(timeout, homePath string) (*relayer.Chain, error) {
	chain, err := ChainInit(timeout, homePath)
	if err != nil {
		return nil, err
	}

	to, err := time.ParseDuration(timeout)
	if err != nil {
		return nil, err
	}

	err = chain.Init(homePath, to, nil, true)
	if err != nil {
		return nil, err
	}

	if chain.KeyExists(chain.Key) {
		logging.Info("deleting old key", chain.Key)

		err = chain.Keybase.Delete(chain.Key)
		if err != nil {
			return nil, err
		}
	}

	_, err = helpers.KeyAddOrRestore(chain, chain.Key, configuration.GetAppConfig().Tendermint.CoinType)
	if err != nil {
		return nil, err
	}

	if err = chain.Start(); err != nil {
		if errors.Is(err, tendermintService.ErrAlreadyStarted) {
			chain.Error(err)

			return nil, err
		}

		return nil, err
	}

	return chain, nil
}

func SetBech32PrefixesAndPStakeWrapAddress(ctx context.Context) (sdkTypes.AccAddress, error) {
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

	// Do not seal the config

	tmAddress, err := casp.GetTendermintAddress(ctx)
	if err != nil {
		return sdkTypes.AccAddress{}, err
	}

	configuration.SetPStakeAddress(tmAddress)

	return tmAddress, nil
}
