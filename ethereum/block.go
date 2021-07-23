package ethereum

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"

	"github.com/Shopify/sarama"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	contracts2 "github.com/persistenceOne/persistenceBridge/ethereum/contracts"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
)

type ethTxToTendermint struct {
	sender common.Address
	msg    sdk.Msg
	txHash string
}

func handleBlock(client *ethclient.Client, ctx *context.Context, block *types.Block, kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec) error {
	var ethTxToTendermintList []ethTxToTendermint
	for _, transaction := range block.Transactions() {
		if transaction.To() != nil {
			var contract contracts2.ContractI
			switch transaction.To().String() {
			case contracts2.LiquidStaking.GetAddress().String():
				contract = &contracts2.LiquidStaking
			case contracts2.TokenWrapper.GetAddress().String():
				contract = &contracts2.TokenWrapper
			default:
			}
			if contract != nil {
				ethTxToTM, err := collectAllEthTxs(client, ctx, transaction, contract)
				if err != nil {
					logging.Error("Failed to process ethereum tx:", transaction.Hash().String())
					return err
				}
				if ethTxToTM.msg != nil {
					ethTxToTendermintList = append(ethTxToTendermintList, ethTxToTM)
				}
			}
		}
	}
	produceToKafka(ethTxToTendermintList, kafkaProducer, protoCodec)
	return nil
}

func collectAllEthTxs(client *ethclient.Client, ctx *context.Context, transaction *types.Transaction, contract contracts2.ContractI) (ethTxToTendermint, error) {
	var ethTxToTM ethTxToTendermint
	receipt, err := client.TransactionReceipt(*ctx, transaction.Hash())
	if err != nil {
		logging.Error("Unable to get receipt of tx:", transaction.Hash().String(), "Error:", err)
		return ethTxToTM, err
	}

	if receipt.Status == 1 {
		logging.Info("Received Ethereum Tx:", transaction.Hash().String())
		method, arguments, err := contract.GetMethodAndArguments(transaction.Data())
		if err != nil {
			return ethTxToTM, fmt.Errorf("unable to get method and arguments of: %s Error: %s", contract.GetName(), err.Error())
		}

		if processFunc, ok := contract.GetSDKMsgAndSender()[method.RawName]; ok {
			msg, sender, err := processFunc(arguments)
			if err != nil {
				return ethTxToTM, fmt.Errorf("failed to process arguments of contract: %s method: %s for TX: %s Error: %s", contract.GetName(), method.RawName, transaction.Hash().String(), err.Error())
			}
			ethTxToTM = ethTxToTendermint{
				sender: sender,
				msg:    msg,
				txHash: transaction.Hash().String(),
			}
			return ethTxToTM, nil
		} else {
			return ethTxToTM, fmt.Errorf("function for method: %s of contract: %s not found", method.RawName, contract.GetName())
		}
	}
	return ethTxToTM, nil
}

func produceToKafka(ethTxToTendermintList []ethTxToTendermint, kafkaProducer *sarama.SyncProducer, protoCodec *codec.ProtoCodec) {
	for _, ethTxToTM := range ethTxToTendermintList {
		msgBytes, err := protoCodec.MarshalInterface(ethTxToTM.msg)
		if err != nil {
			logging.Fatal(err)
		}
		producer := ""
		switch ethTxToTM.msg.Type() {
		case bankTypes.TypeMsgSend:
			producer = utils.MsgSend
		case stakingTypes.TypeMsgDelegate:
			producer = utils.MsgDelegate
		case stakingTypes.TypeMsgUndelegate:
			producer = utils.EthUnbond
		default:
			logging.Fatal("unknown msg type [ETH Listener]")
		}
		logging.Info("Adding to [ETH Listener] kafka producer:", producer, "of txHash:", ethTxToTM.txHash, "msgType:", ethTxToTM.msg.Type(), "sender:", ethTxToTM.sender.String())
		err = utils.ProducerDeliverMessage(msgBytes, producer, *kafkaProducer)
		if err != nil {
			logging.Fatal("Failed to add msg to kafka queue [ETH Listener], producer:", producer, "txHash:", ethTxToTM.txHash, "msg:", ethTxToTM.msg.String(), "sender:", ethTxToTM.sender.String(), "error:", err)
		}
	}
}
