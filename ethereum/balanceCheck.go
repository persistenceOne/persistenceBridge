package ethereum

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
	"math/big"
)

func BalanceCheck(currentHeight uint64, client *ethclient.Client) {
	if configuration.GetAppConfig().Ethereum.BalanceCheckPeriod != 0 && currentHeight%configuration.GetAppConfig().Ethereum.BalanceCheckPeriod == 0 {
		balance, err := client.BalanceAt(context.Background(), configuration.GetAppConfig().Ethereum.GetBridgeAdminAddress(), nil)
		if err != nil {
			logging.Error("Unable to fetch eth bridge admin balance")
		}
		ethAlertAmount := big.NewInt(0).Mul(big.NewInt(configuration.GetAppConfig().Ethereum.AlertAmount), big.NewInt(1000000000))
		if balance.Cmp(ethAlertAmount) <= 0 {
			logging.Warn("Ethereum bridge admin address", configuration.GetAppConfig().Ethereum.GetBridgeAdminAddress().String(), "balance has fallen below  ")
		}
	}
}
