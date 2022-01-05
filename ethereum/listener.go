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

	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/application/shutdown"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

const finality = 12

func StartListening(client *ethclient.Client, database *badger.DB, sleepDuration time.Duration, brokers []string, protoCodec *codec.ProtoCodec) {
	ctx := context.Background()
	kafkaProducer := utils.NewProducer(brokers, utils.SaramaConfig())

	defer func(kafkaProducer sarama.SyncProducer) {
		err := kafkaProducer.Close()
		if err != nil {
			logging.Error(err)
		}
	}(kafkaProducer)

	for {
		if shutdown.GetBridgeStopSignal() {
			if shutdown.GetKafkaConsumerClosed() {
				logging.Info("Stopping Ethereum Listener!!!")

				shutdown.SetETHStopped(true)

				return
			}

			time.Sleep(1 * time.Second)

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
			logging.Error("Stopping Ethereum Listener, unable to get status, Error:", err)

			shutdown.SetETHStopped(true)

			return
		}

		if (latestEthHeight - uint64(ethStatus.LastCheckHeight)) > finality {
			processHeight := big.NewInt(ethStatus.LastCheckHeight + 1)

			logging.Info("Ethereum Block:", processHeight)

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
				logging.Error("Stopping Ethereum Listener, unable to set (DB) status to", processHeight, "Error:", err)

				shutdown.SetETHStopped(true)

				return
			}

			err = onNewBlock(ctx, latestEthHeight, client, kafkaProducer, database)
			if err != nil {
				logging.Error("Stopping Ethereum Listener, onNewBlock error:", err)

				shutdown.SetETHStopped(true)

				return
			}
		}

		time.Sleep(sleepDuration)
	}
}
