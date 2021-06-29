package tendermint

import (
	"context"
	"fmt"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/application/shutdown"
	"log"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/relayer/relayer"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
)

func StartListening(initClientCtx client.Context, chain *relayer.Chain, kafkaState utils.KafkaState, protoCodec *codec.ProtoCodec, sleepDuration time.Duration) {
	ctx := context.Background()

	for {
		if shutdown.GetBridgeStopSignal() {
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
				log.Println(err)
				time.Sleep(sleepDuration)
				continue
			}

			err = handleTxSearchResult(initClientCtx, txSearchResult, kafkaState, protoCodec)
			if err != nil {
				panic(err)
			}

			err = db.SetCosmosStatus(processHeight)
			if err != nil {
				panic(err)
			}
		}

		err = onNewBlock(ctx, chain, kafkaState, protoCodec)
		if err != nil {
			panic(err)
		}
		time.Sleep(sleepDuration)
	}
}
