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
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/relayer/relayer"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

func onNewBlock(ctx context.Context, clientCtx client.Context, chain *relayer.Chain, kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec) error {
	return db.IterateOutgoingTmTx(func(key []byte, value []byte) error {
		var tmTx db.OutgoingTendermintTransaction
		err := json.Unmarshal(value, &tmTx)
		if err != nil {
			return fmt.Errorf("failed to unmarshal OutgoingTendermintTransaction %s [TM onNewBlock]: %s", string(key), err.Error())
		}
		txHashBytes, err := hex.DecodeString(tmTx.TxHash)
		if err != nil {
			return fmt.Errorf("invalid tx hash %s [TM onNewBlock]: %s", tmTx.TxHash, err.Error())
		}
		txResult, err := chain.Client.Tx(ctx, txHashBytes, true)
		if err != nil {
			if err.Error() == fmt.Sprintf("RPC error -32603 - Internal error: tx (%s) not found", tmTx.TxHash) {
				logging.Info("Tendermint tx still pending:", tmTx.TxHash)
				return nil
			}
			logging.Error(fmt.Sprintf("Tendermint tx hash %s search failed [TM onNewBlock]: %s", tmTx.TxHash, err.Error()))
			return err
		} else {
			if txResult.TxResult.GetCode() != 0 {
				logging.Error(fmt.Sprintf("Broadcast tendermint tx %s (block: %d) failed, Code: %d, Log: %s", tmTx.TxHash, txResult.Height, txResult.TxResult.GetCode(), txResult.TxResult.Log))
				txInterface, err := clientCtx.TxConfig.TxDecoder()(txResult.Tx)
				if err != nil {
					return err
				}

				transaction, ok := txInterface.(signing.Tx)
				if !ok {
					return fmt.Errorf("unable to parse transaction into signing.Tx [TM onNewBlock]")
				}
				for _, msg := range transaction.GetMsgs() {
					if sdkTypes.MsgTypeURL(msg) != constants.MsgWithdrawDelegatorRewardTypeUrl {
						msgBytes, err := protoCodec.MarshalInterface(msg)
						if err != nil {
							return fmt.Errorf("failed to generate msgBytes [TM onNewBlock]: %s", err.Error())
						}
						err = utils.ProducerDeliverMessage(msgBytes, utils.RetryTendermint, *kafkaProducer)
						if err != nil {
							return fmt.Errorf("failed to add messages of %s to kafka queue [TM onNewBlock] RetryTendermint: %s", msg.String(), err.Error())
						}
					}
				}
			} else {
				logging.Info("Broadcast tendermint tx successful. Hash:", tmTx.TxHash, "Block:", txResult.Height)
			}
			return db.DeleteOutgoingTendermintTx(tmTx.TxHash)
		}
	})
}
