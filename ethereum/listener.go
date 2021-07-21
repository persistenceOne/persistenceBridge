package ethereum

import (
	"context"
	"github.com/Shopify/sarama"
	"math/big"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/application/shutdown"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

func StartListening(client *ethclient.Client, sleepDuration time.Duration, brokers []string, protoCodec *codec.ProtoCodec) {
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

		ethStatus, err := db.GetEthereumStatus()
		if err != nil {
			panic(err)
		}

		if (latestEthHeight - uint64(ethStatus.LastCheckHeight)) > 12 {
			processHeight := big.NewInt(ethStatus.LastCheckHeight + 1)
			logging.Info("Ethereum Block:", processHeight)

			block, err := client.BlockByNumber(ctx, processHeight)
			if err != nil {
				logging.Error("Unable to fetch ethereum block:", processHeight, "Error:", err)
				time.Sleep(sleepDuration)
				continue
			}

			err = handleBlock(client, &ctx, block, &kafkaProducer, protoCodec)
			if err != nil {
				logging.Fatal("Ethereum listener unable to handleBlock:", processHeight, "Error:", err)
			}

			err = db.SetEthereumStatus(processHeight.Int64())
			if err != nil {
				logging.Fatal("setting ethereum status:", err)
			}

			err = onNewBlock(ctx, latestEthHeight, client, &kafkaProducer)
			if err != nil {
				logging.Error("Stopping ethereum Listener, onNewBlock error:", err)
				shutdown.SetETHStopped(true)
				return
			}
		}
		time.Sleep(sleepDuration)
	}
}
