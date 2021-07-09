package tendermint

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/application/shutdown"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"log"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/relayer/relayer"
)

func StartListening(initClientCtx client.Context, chain *relayer.Chain, brokers []string, protoCodec *codec.ProtoCodec, sleepDuration time.Duration) {
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
		if shutdown.GetBridgeStopSignal() && shutdown.GetKafkaConsumerClosed() {
			log.Println("Stopping Tendermint Listener!!!")
			shutdown.SetTMStopped(true)
			return
		}

		abciInfo, err := chain.Client.ABCIInfo(ctx)
		if err != nil {
			log.Printf("Error while fetching tendermint abci info: %s\n", err.Error())
			time.Sleep(sleepDuration)
			continue
		}

		cosmosStatus, err := db.GetCosmosStatus()
		if err != nil {
			panic(err)
		}

		if abciInfo.Response.LastBlockHeight > cosmosStatus.LastCheckHeight {
			processHeight := cosmosStatus.LastCheckHeight + 1
			log.Printf("TM: %d\n", processHeight)

			txSearchResult, err := chain.Client.TxSearch(ctx, fmt.Sprintf("tx.height=%d", processHeight), true, nil, nil, "asc")
			if err != nil {
				log.Printf("ERROR getting TM height %d: %s\n", processHeight, err.Error())
				time.Sleep(sleepDuration)
				continue
			}

			err = handleTxSearchResult(initClientCtx, txSearchResult, &kafkaProducer, protoCodec)
			if err != nil {
				panic(err)
			}

			err = db.SetCosmosStatus(processHeight)
			if err != nil {
				panic(err)
			}
		}

		err = onNewBlock(ctx, initClientCtx, chain, &kafkaProducer, protoCodec)
		if err != nil {
			panic(err)
		}
		time.Sleep(sleepDuration)
	}
}
