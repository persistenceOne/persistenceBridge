/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

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
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/relayer/relayer"

	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

func onNewBlock(ctx context.Context, clientCtx *client.Context, chain *relayer.Chain, kafkaProducer sarama.SyncProducer, protoCodec *codec.ProtoCodec) error {
	return db.IterateOutgoingTmTx(func(key []byte, value []byte) error {
		var tmTx db.OutgoingTendermintTransaction

		err := json.Unmarshal(value, &tmTx)
		if err != nil {
			return fmt.Errorf("%w %s [TM onNewBlock]: %s", ErrUnmarshalOutgoingTransaction, string(key), err.Error())
		}

		txHashBytes, err := hex.DecodeString(tmTx.TxHash)
		if err != nil {
			return fmt.Errorf("%w %s [TM onNewBlock]: %s", ErrInvalidTxHash, tmTx.TxHash, err.Error())
		}

		txResult, err := chain.Client.Tx(ctx, txHashBytes, true)
		if err != nil {
			if err.Error() == fmt.Sprintf("RPC bridgeErr -32603 - Internal bridgeErr: tx (%s) not found", tmTx.TxHash) {
				logging.Info("Tendermint tx still pending:", tmTx.TxHash)

				return nil
			}

			logging.Error(fmt.Sprintf("Tendermint tx hash %s search failed [TM onNewBlock]: %s", tmTx.TxHash, err.Error()))

			return err
		}

		if txResult.TxResult.GetCode() != 0 {
			logging.Error(fmt.Sprintf("Broadcast tendermint tx %s (block: %d) failed, Code: %d, Log: %s", tmTx.TxHash, txResult.Height, txResult.TxResult.GetCode(), txResult.TxResult.Log))

			txInterface, err := clientCtx.TxConfig.TxDecoder()(txResult.Tx)
			if err != nil {
				return err
			}

			transaction, ok := txInterface.(signing.Tx)
			if !ok {
				return fmt.Errorf("%w [TM onNewBlock]", ErrParseTransaction)
			}

			for _, msg := range transaction.GetMsgs() {
				if msg.Type() != distributionTypes.TypeMsgWithdrawDelegatorReward {
					msgBytes, err := protoCodec.MarshalInterface(msg)
					if err != nil {
						return fmt.Errorf("%w [TM onNewBlock]: %s", ErrTransactionMessageGeneration, err.Error())
					}

					err = utils.ProducerDeliverMessage(msgBytes, utils.RetryTendermint, kafkaProducer)
					if err != nil {
						return fmt.Errorf("%w [TM onNewBlock] RetryTendermint: message %q, error %s",
							ErrAddToKafkaQueue, msg.String(), err.Error())
					}
				}
			}
		} else {
			logging.Info("Broadcast tendermint tx successful. Hash:", tmTx.TxHash, "Block:", txResult.Height)
		}

		return db.DeleteOutgoingTendermintTx(tmTx.TxHash)
	})
}
