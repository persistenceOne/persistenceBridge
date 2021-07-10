package ethereum

import (
	"context"
	"github.com/Shopify/sarama"
	"log"
	"math/big"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/application/shutdown"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
)

func StartListening(client *ethclient.Client, sleepDuration time.Duration, brokers []string, protoCodec *codec.ProtoCodec) {
	ctx := context.Background()
	kafkaProducer := utils.NewProducer(brokers, utils.SaramaConfig())
	defer func(kafkaProducer sarama.SyncProducer) {
		err := kafkaProducer.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}(kafkaProducer)

	for {
		if shutdown.GetBridgeStopSignal() {
			if shutdown.GetKafkaConsumerClosed() {
				log.Println("Stopping Ethereum Listener!!!")
				shutdown.SetETHStopped(true)
				return
			}
			time.Sleep(1 * time.Second)
			continue
		}

		latestEthHeight, err := client.BlockNumber(ctx)
		if err != nil {
			log.Printf("Error while fetching latest block height: %s\n", err.Error())
			time.Sleep(sleepDuration)
			continue
		}

		ethStatus, err := db.GetEthereumStatus()
		if err != nil {
			panic(err)
		}

		if (latestEthHeight - uint64(ethStatus.LastCheckHeight)) > 12 {
			processHeight := big.NewInt(ethStatus.LastCheckHeight + 1)
			log.Printf("ETH: %d\n", processHeight)

			block, err := client.BlockByNumber(ctx, processHeight)
			if err != nil {
				log.Printf("ERROR getting ETH height %d: %s\n", processHeight, err.Error())
				time.Sleep(sleepDuration)
				continue
			}

			err = handleBlock(client, &ctx, block, &kafkaProducer, protoCodec)
			if err != nil {
				panic(err)
			}

			err = db.SetEthereumStatus(processHeight.Int64())
			if err != nil {
				panic(err)
			}

			err = onNewBlock(ctx, latestEthHeight, client, &kafkaProducer)
			if err != nil {
				panic(err)
			}
		}
		time.Sleep(sleepDuration)
	}
}
