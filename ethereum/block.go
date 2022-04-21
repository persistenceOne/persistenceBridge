/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package ethereum

import (
	"context"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	"github.com/persistenceOne/persistenceBridge/application/db"
	"github.com/persistenceOne/persistenceBridge/ethereum/abi/liquidStaking"
	"github.com/persistenceOne/persistenceBridge/ethereum/abi/tokenWrapper"
	"github.com/persistenceOne/persistenceBridge/ethereum/contracts"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

const (
	EventStakeTokens     = "StakeTokens"
	EventUnstakeTokens   = "UnstakeTokens"
	EventWithdrawUTokens = "WithdrawUTokens"
)

func handleBlock(client *ethclient.Client, ctx *context.Context, block *types.Block, kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec) error {

	contracts.LiquidStaking.GetAddress()
	logs, err := client.FilterLogs(*ctx, ethereum.FilterQuery{
		//BlockHash: &ab,
		FromBlock: block.Number(),
		ToBlock:   block.Number(),
		Addresses: []common.Address{contracts.LiquidStaking.GetAddress(), contracts.TokenWrapper.GetAddress()},
		Topics: [][]common.Hash{{contracts.LiquidStaking.GetABI().Events[EventStakeTokens].ID,
			contracts.LiquidStaking.GetABI().Events[EventUnstakeTokens].ID,
			contracts.TokenWrapper.GetABI().Events[EventWithdrawUTokens].ID}},
	})
	if err != nil {
		return err
	}

	for _, vLog := range logs {
		// check if log is removed
		if vLog.Removed {
			// Test this out somehow
			// logging.Warn(fmt.Sprintf("Ethereum Reader: log was removed on Block %s, txhash: %s", block.Number(), vLog.TxHash))
			continue
		}
		// check if txn was successful
		receipt, err := client.TransactionReceipt(*ctx, vLog.TxHash)
		if err != nil {
			logging.Error("Unable to get receipt of tx:", vLog.TxHash.String(), "Error:", err)
			return err
		}
		if !(receipt.Status == types.ReceiptStatusSuccessful) {
			continue
		}

		var contract contracts.ContractI
		var methodRaw string
		var args []interface{}
		switch vLog.Topics[0] {
		case contracts.LiquidStaking.GetABI().Events[EventStakeTokens].ID:
			stakeEvent, err := ParseLiquidStakingStakeTokensEvent(vLog)
			if err != nil {
				logging.Error("Failed to process ethereum tx:", vLog.TxHash.String(), "with err: ", err.Error())
				return err
			}
			methodRaw = constants.LiquidStakingStake
			contract = &contracts.LiquidStaking
			args = []interface{}{stakeEvent.AccountAddress, stakeEvent.FinalTokens}

		case contracts.LiquidStaking.GetABI().Events[EventUnstakeTokens].ID:
			unstakeEvent, err := ParseLiquidStakingUnstakeTokensEvent(vLog)
			if err != nil {
				logging.Error("Failed to process ethereum tx:", vLog.TxHash.String(), "with err: ", err.Error())
				return err
			}
			methodRaw = constants.LiquidStakingUnStake
			contract = &contracts.LiquidStaking
			args = []interface{}{unstakeEvent.AccountAddress, unstakeEvent.FinalTokens}
		case contracts.TokenWrapper.GetABI().Events[EventWithdrawUTokens].ID:
			withdrawUTokenEvent, err := ParseTokenWrapperWithdrawUTokensEvent(vLog)
			if err != nil {
				logging.Error("Failed to process ethereum tx:", vLog.TxHash.String(), "with err: ", err.Error())
				return err
			}
			methodRaw = constants.TokenWrapperWithdrawUTokens
			contract = &contracts.TokenWrapper
			args = []interface{}{withdrawUTokenEvent.AccountAddress, withdrawUTokenEvent.FinalTokens, withdrawUTokenEvent.ToChainAddress}
		default:
		}
		err = collectEthTx(protoCodec, vLog.TxHash, methodRaw, contract, args)
		if err != nil {
			logging.Error("Failed to process ethereum tx:", vLog.TxHash.String())
			return err
		}
	}

	produceToKafka(kafkaProducer)
	return nil
}

func collectEthTx(protoCodec *codec.ProtoCodec, txHash common.Hash, methodRawName string, contract contracts.ContractI, arguments []interface{}) error {

	if processFunc, ok := contract.GetSDKMsgAndSender()[methodRawName]; ok {
		msg, sender, err := processFunc(arguments)
		if err != nil {
			return fmt.Errorf("failed to process arguments of contract: %s method: %s for TX: %s Error: %s", contract.GetName(), methodRawName, txHash.String(), err.Error())
		}
		// Do not check for EthereumTxToKafka exists.
		if !db.CheckIncomingEthereumTxExists(txHash) {
			msgBytes, err := protoCodec.MarshalInterface(msg)
			if err != nil {
				return err
			}
			err = db.AddIncomingEthereumTx(db.IncomingEthereumTx{
				TxHash:   txHash,
				MsgBytes: msgBytes,
				Sender:   sender,
				MsgType:  sdkTypes.MsgTypeURL(msg),
			})
			if err != nil {
				return err
			}
			err = db.AddEthereumTxToKafka(db.EthereumTxToKafka{
				TxHash: txHash,
			})
			if err != nil {
				return fmt.Errorf("added to IncomingEthereumTx but NOT to EthereumTxToKafka failed. Tx won't be added to kafka: %v", err)
			}
		}
	}

	return nil
}

func produceToKafka(kafkaProducer *sarama.SyncProducer) {
	ethTxToKafkaList, err := db.GetAllEthereumTxToKafka()
	if err != nil {
		logging.Fatal(err)
	}
	for _, tx := range ethTxToKafkaList {
		ethTxToTM, err := db.GetIncomingEthereumTx(tx.TxHash)
		if err != nil {
			logging.Fatal(err)
		}
		producer := ""

		switch ethTxToTM.MsgType {
		case constants.MsgSendTypeUrl:
			producer = utils.MsgSend
		case constants.MsgDelegateTypeUrl:
			producer = utils.MsgDelegate
		case constants.MsgUndelegateTypeUrl:
			producer = utils.EthUnbond
		default:
			logging.Fatal("unknown msg type [ETH Listener]: ", ethTxToTM.MsgType)
		}
		logging.Info("[ETH Listener] Adding to kafka producer:", producer, "of txHash:", ethTxToTM.TxHash.String(), "msgType:", ethTxToTM.MsgType, "sender:", ethTxToTM.Sender.String())
		err = utils.ProducerDeliverMessage(ethTxToTM.MsgBytes, producer, *kafkaProducer)
		if err != nil {
			logging.Fatal("Failed to add msg to kafka queue [ETH Listener], producer:", producer, "txHash:", ethTxToTM.TxHash.String(), "sender:", ethTxToTM.Sender.String(), "error:", err)
		}
		err = db.DeleteEthereumTxToKafka(ethTxToTM.TxHash)
		if err != nil {
			logging.Fatal(err)
		}
	}
}

//Helper functions

func ParseLiquidStakingStakeTokensEvent(vLog types.Log) (liquidStaking.LiquidStakingStakeTokens, error) {
	var event liquidStaking.LiquidStakingStakeTokens
	err := contracts.LiquidStaking.GetABI().UnpackIntoInterface(&event, EventStakeTokens, vLog.Data)
	if err != nil {
		return event, errors.New("failed to Unpack event StakeTokens")
	}
	event.Raw = vLog
	event.AccountAddress = common.HexToAddress(vLog.Topics[1].Hex())
	event.Tokens = vLog.Topics[2].Big()
	event.FinalTokens = vLog.Topics[3].Big()

	return event, nil
}
func ParseLiquidStakingUnstakeTokensEvent(vLog types.Log) (liquidStaking.LiquidStakingUnstakeTokens, error) {
	var event liquidStaking.LiquidStakingUnstakeTokens
	err := contracts.LiquidStaking.GetABI().UnpackIntoInterface(&event, EventUnstakeTokens, vLog.Data)
	if err != nil {
		return event, errors.New("failed to Unpack event UnstakeTokens")
	}
	event.Raw = vLog
	event.AccountAddress = common.HexToAddress(vLog.Topics[1].Hex())
	event.Tokens = vLog.Topics[2].Big()
	event.FinalTokens = vLog.Topics[3].Big()
	return event, nil
}
func ParseTokenWrapperWithdrawUTokensEvent(vLog types.Log) (tokenWrapper.TokenWrapperWithdrawUTokens, error) {
	var event tokenWrapper.TokenWrapperWithdrawUTokens
	err := contracts.TokenWrapper.GetABI().UnpackIntoInterface(&event, EventWithdrawUTokens, vLog.Data)
	if err != nil {
		return event, errors.New("failed to Unpack event WithdrawUTokens")
	}
	event.Raw = vLog
	event.AccountAddress = common.HexToAddress(vLog.Topics[1].Hex())
	event.Tokens = vLog.Topics[2].Big()
	return event, nil
}
