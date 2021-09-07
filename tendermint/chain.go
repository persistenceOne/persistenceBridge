package tendermint

import (
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"time"

	"github.com/cosmos/relayer/helpers"
	"github.com/cosmos/relayer/relayer"
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

	//118 is not being used anywhere
	ko, err := helpers.KeyAddOrRestore(chain, chain.Key, uint32(118))
	if err != nil {
		return chain, err
	}

	logging.Info("Relayer Chain Keys added [NOT TO BE USED]:", ko.Address)

	if err = chain.Start(); err != nil {
		if err != tendermintService.ErrAlreadyStarted {
			chain.Error(err)
			return chain, err
		}
	}
	return chain, nil
}
