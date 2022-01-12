/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package ethereum

import (
	"context"
	"math/big"
	"time"

	"github.com/Shopify/sarama"
	"github.com/dgraph-io/badger/v3"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/application/shutdown"
	"github.com/persistenceOne/persistenceBridge/ethereum/contracts"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

func StartListening(client *ethclient.Client, database *badger.DB, sleepDuration time.Duration, brokers []string, protoCodec *codec.ProtoCodec) {
	ctx := context.Background()
	kafkaProducer := utils.NewProducer(brokers, utils.SaramaConfig())

	defer func(kafkaProducer sarama.SyncProducer) {
		err := kafkaProducer.Close()
		if err != nil {
			logging.Error(err)
		}
	}(kafkaProducer)

	// Need to set it here because configuration isn't initialized when contract objects are created
	contracts.LiquidStaking.SetAddress(configuration.GetAppConfig().Ethereum.LiquidStakingAddress)
	contracts.TokenWrapper.SetAddress(configuration.GetAppConfig().Ethereum.TokenWrapperAddress)

	ethAlertAmount = big.NewInt(0).Mul(big.NewInt(configuration.GetAppConfig().Ethereum.AlertAmount), big.NewInt(1000000000))

	for {
		if shutdown.GetBridgeStopSignal() {
			if shutdown.GetKafkaConsumerClosed() {
				logging.Info("Stopping Ethereum Listener!!!")

				shutdown.SetETHStopped(true)

				return
			}

			// fixme: use timer.Ticker
			time.Sleep(1 * time.Second) // thread is put to sleep to prevent 100% CPU usage

			continue
		}

		latestEthHeight, err := client.BlockNumber(ctx)
		if err != nil {
			logging.Error("Unable to fetch ethereum latest block height:", err)

			time.Sleep(sleepDuration)

			continue
		}

		ethStatus, err := db.GetEthereumStatus(database)
		if err != nil {
			logging.Fatal("Stopping Ethereum Listener, unable to get status, Error:", err)
		}

		if ethStatus.LastCheckHeight < 0 {
			logging.Fatal("Stopping Ethereum Listener, eth status height is less than 0:", ethStatus.LastCheckHeight)
		}

		if (latestEthHeight - uint64(ethStatus.LastCheckHeight)) > constants.EthereumBlockConfirmations {
			processHeight := big.NewInt(ethStatus.LastCheckHeight + 1)

			logging.Info("Ethereum Block:", processHeight)

			BalanceCheck(processHeight.Uint64(), client)

			block, err := client.BlockByNumber(ctx, processHeight)
			if err != nil {
				logging.Error("Unable to fetch ethereum block:", processHeight, "Error:", err)

				time.Sleep(sleepDuration)

				continue
			}

			err = handleBlock(ctx, client, database, block, kafkaProducer, protoCodec)
			if err != nil {
				logging.Error("Unable to fetch handle ethereum block:", processHeight, "Error:", err)

				time.Sleep(sleepDuration)

				continue
			}

			err = db.SetEthereumStatus(database, processHeight.Int64())
			if err != nil {
				logging.Fatal("Stopping Ethereum Listener, unable to set (DB) status to", processHeight, "Error:", err)
			}

			err = onNewBlock(ctx, latestEthHeight, client, kafkaProducer, protoCodec, database)
			if err != nil {
				logging.Fatal("Stopping Ethereum Listener, onNewBlock error:", err)
			}
		}

		time.Sleep(sleepDuration)
	}
}
