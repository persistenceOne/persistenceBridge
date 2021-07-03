package ethereum

import (
	"context"
	"log"
	"math/big"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/persistenceOne/persistenceBridge/application"
	"github.com/persistenceOne/persistenceBridge/application/shutdown"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
)

func StartListening(client *ethclient.Client, sleepDuration time.Duration, kafkaState utils.KafkaState, protoCodec *codec.ProtoCodec) {
	ctx := context.Background()

	for {
		if shutdown.GetBridgeStopSignal() {
			log.Println("Stopping Ethereum Listener!!!")
			shutdown.SetETHStopped(true)
			return
		}

		latestEthHeight, err := client.BlockNumber(ctx)
		if err != nil {
			log.Printf("Error while fetching latest block height: %s\n", err.Error())
			time.Sleep(sleepDuration)
			continue
		}

		ethStatus, err := application.GetEthereumStatus()
		if err != nil {
			panic(err)
		}

		if latestEthHeight > uint64(ethStatus.LastCheckHeight) {
			processHeight := big.NewInt(ethStatus.LastCheckHeight + 1)
			log.Printf("ETH: %d\n", processHeight)

			block, err := client.BlockByNumber(ctx, processHeight)
			if err != nil {
				log.Printf("ERROR getting ETH height %d: %s\n", processHeight, err.Error())
				time.Sleep(sleepDuration)
				continue
			}

			err = handleBlock(client, &ctx, block, kafkaState, protoCodec)
			if err != nil {
				panic(err)
			}

			err = application.SetEthereumStatus(processHeight.Int64())
			if err != nil {
				panic(err)
			}
		}
		time.Sleep(sleepDuration)
	}
}
