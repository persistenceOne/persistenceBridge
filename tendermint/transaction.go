package tendermint

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	goEthCommon "github.com/ethereum/go-ethereum/common"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/application/outgoingTx"
	"github.com/persistenceOne/persistenceBridge/application/shutdown"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
	tmCoreTypes "github.com/tendermint/tendermint/rpc/core/types"
)

func handleTxSearchResult(clientCtx client.Context, resultTxs []*tmCoreTypes.ResultTx, kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec) error {
	for _, transaction := range resultTxs {
		err := collectAllWrapAndRevertTxs(clientCtx, transaction)
		if err != nil {
			logging.Error("Failed to process tendermint transaction:", transaction.Hash.String())
			return err
		}
	}
	wrapOrRevert(kafkaProducer, protoCodec)
	return nil
}

func collectAllWrapAndRevertTxs(clientCtx client.Context, txQueryResult *tmCoreTypes.ResultTx) error {
	if txQueryResult.TxResult.GetCode() == 0 {
		// Should be used if txQueryResult.Tx is string
		//decodedTx, err := base64.StdEncoding.DecodeString(txQueryResult.Tx)
		//if err != nil {
		//	return txMsgs, err
		//}

		txInterface, err := clientCtx.TxConfig.TxDecoder()(txQueryResult.Tx)
		if err != nil {
			return err
		}

		transaction, ok := txInterface.(signing.Tx)
		if !ok {
			return err
		}

		memo := strings.TrimSpace(transaction.GetMemo())

		for i, msg := range transaction.GetMsgs() {
			switch txMsg := msg.(type) {
			case *banktypes.MsgSend:
				if txMsg.ToAddress == configuration.GetAppConfig().Tendermint.GetPStakeAddress() {
					if memo != "DO_NOT_REVERT" {
						for _, coin := range txMsg.Amount {
							produced := db.CheckTendermintIncomingTxProduced(txQueryResult.Hash, uint(i), coin.Denom)
							if !produced {
								err = db.AddToPendingTendermintIncomingTx(db.TendermintIncomingTx{
									TxHash:          txQueryResult.Hash,
									MsgIndex:        uint(i),
									ProducedToKafka: false,
									Denom:           coin.Denom,
									FromAddress:     txMsg.FromAddress,
									Amount:          coin.Amount,
									Memo:            memo,
								})
								if err != nil {
									return err
								}
							}
						}
					} else {
						logging.Info("Deposited to wrapping address, TxHash:", txQueryResult.Hash.String(), "amount:", txMsg.Amount.String())
					}
				}
			default:
			}
		}
	}
	return nil
}

func wrapOrRevert(kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec) {
	txs, err := db.GetProduceToKafkaTendermintIncomingTxs()
	if err != nil {
		logging.ShutdownWithError(err)
		shutdown.SetBridgeStopSignal(true)
	}
	for _, tx := range txs {
		validEthMemo := goEthCommon.IsHexAddress(tx.Memo)
		var ethAddress goEthCommon.Address
		if validEthMemo {
			ethAddress = goEthCommon.HexToAddress(tx.Memo)
		}

		// TODO tx.Amount.GT required??
		if tx.Denom == configuration.GetAppConfig().Tendermint.PStakeDenom && validEthMemo && tx.Amount.GT(sdk.ZeroInt()) {
			logging.Info("Tendermint Wrap Tx:", tx.TxHash, "Msg Index:", tx.MsgIndex, "Amount:", tx.Amount.String())
			ethTxMsg := outgoingTx.WrapTokenMsg{
				Address: ethAddress,
				Amount:  tx.Amount.BigInt(),
			}
			msgBytes, err := json.Marshal(ethTxMsg)
			if err != nil {
				panic(err)
			}
			logging.Info("Adding wrap token msg to kafka producer ToEth, from:", tx.FromAddress, "to:", ethAddress.String(), "amount:", tx.Amount.String())
			err = utils.ProducerDeliverMessage(msgBytes, utils.ToEth, *kafkaProducer)
			if err != nil {
				logging.ShutdownWithError(fmt.Errorf("failed to add msg to kafka queue ToEth [TM Listener]: %s", err.Error()))
				shutdown.SetBridgeStopSignal(true)
			}
		} else {
			revertToken := sdk.NewCoin(tx.Denom, tx.Amount)
			logging.Info("Reverting coins,TxHash:", tx.TxHash, "Msg Index:", tx.MsgIndex, "Coin:", revertToken.String())
			revertCoins(tx.FromAddress, sdk.NewCoins(revertToken), kafkaProducer, protoCodec)
		}
		err = db.SetTendermintIncomingTxProduced(tx)
		if err != nil {
			logging.ShutdownWithError(err)
			shutdown.SetBridgeStopSignal(true)
		}
	}
}

func revertCoins(toAddress string, coins sdk.Coins, kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec) {
	msg := &banktypes.MsgSend{
		FromAddress: configuration.GetAppConfig().Tendermint.GetPStakeAddress(),
		ToAddress:   toAddress,
		Amount:      coins,
	}
	msgBytes, err := protoCodec.MarshalInterface(msg)
	if err != nil {
		logging.ShutdownWithError(err)
		shutdown.SetBridgeStopSignal(true)
	}
	logging.Info("REVERT: adding send coin msg to kafka producer MsgSend, to:", toAddress, "amount:", coins.String())
	err = utils.ProducerDeliverMessage(msgBytes, utils.MsgSend, *kafkaProducer)
	if err != nil {
		logging.ShutdownWithError(fmt.Errorf("failed to add msg to kafka queue ToEth [TM Listener REVERT]: %s", err.Error()))
		shutdown.SetBridgeStopSignal(true)
	}
}
