package tendermint

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/relayer/relayer"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"log"
)

func onNewBlock(ctx context.Context, clientCtx client.Context, chain *relayer.Chain, kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec) error {
	return db.IterateTmTx(func(key []byte, value []byte) error {
		var tmTx db.TendermintBroadcastedTransaction
		err := json.Unmarshal(value, &tmTx)
		if err != nil {
			log.Fatalln("Failed to unmarshal TendermintBroadcastedTransaction : ", err)
		}
		txHashBytes, err := hex.DecodeString(tmTx.TxHash)
		if err != nil {
			log.Fatalln("invalid tx hash : ", err)
		}
		txResult, err := chain.Client.Tx(ctx, txHashBytes, true)
		if err != nil {
			if err.Error() == fmt.Sprintf("RPC error -32603 - Internal error: tx (%s) not found", tmTx.TxHash) {
				log.Printf("TM Tx %s still pending\n", tmTx.TxHash)
				return nil
			}
			log.Printf("tm tx hash search failed: %s\n", err)
			return err
		} else {
			if txResult.TxResult.GetCode() != 0 {
				log.Printf("Broadcasted tendermint tx %s (block: %d) failed, code: %d, log: %s\n", tmTx.TxHash, txResult.Height, txResult.TxResult.GetCode(), txResult.TxResult.Log)
				txInterface, err := clientCtx.TxConfig.TxDecoder()(txResult.Tx)
				if err != nil {
					log.Fatalln(err.Error())
				}

				transaction, ok := txInterface.(signing.Tx)
				if !ok {
					log.Fatalln("Unable to parse transaction into signing.Tx")
				}
				for _, msg := range transaction.GetMsgs() {
					msgBytes, err := protoCodec.MarshalInterface(msg)
					if err != nil {
						log.Fatalln("Failed to generate msgBytes: ", err)
					}
					err = utils.ProducerDeliverMessage(msgBytes, utils.RetryTendermint, *kafkaProducer)
					if err != nil {
						log.Fatalf("Failed to add messages of %s to kafka queue RetryTendermint: %s [TM onNewBlock]\n", tmTx.TxHash, err.Error())
					}
				}
			} else {
				log.Printf("Broadcasted tendermint tx %s (block: %d) success.\n", tmTx.TxHash, txResult.Height)
			}
			return db.DeleteTendermintTx(tmTx.TxHash)
		}

		return nil
	})
}
