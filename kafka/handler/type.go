package handler

import (
	"encoding/json"
	"errors"
	"github.com/Shopify/sarama"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/relayer/relayer"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/persistenceOne/persistenceCore/kafka/runconfig"
	"github.com/persistenceOne/persistenceCore/kafka/utils"
	"github.com/persistenceOne/persistenceCore/pStake/constants"
	"github.com/persistenceOne/persistenceCore/pStake/ethereum"
	"log"
)

type MsgHandler struct {
	KafkaConfig runconfig.KafkaConfig
	ProtoCodec  *codec.ProtoCodec
	Chain       *relayer.Chain
	EthClient   *ethclient.Client
}

var _ sarama.ConsumerGroupHandler = MsgHandler{}

func (m MsgHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (m MsgHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (m MsgHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	switch claim.Topic() {
	case utils.ToEth:
		err := m.HandleTopicMsgs(session, claim, m.KafkaConfig.ToEth.BatchSize, SendBatchToEth)
		if err != nil {
			log.Printf("failed batch and handle for topic: %v with error %v", utils.ToEth, err)
			return err
		}
	case utils.ToTendermint:
		err := m.HandleTopicMsgs(session, claim, m.KafkaConfig.ToTendermint.BatchSize, SendBatchToTendermint)
		if err != nil {
			log.Printf("failed batch and handle for topic: %v with error %v", utils.ToTendermint, err)
			return err
		}
	case utils.EthUnbond:
		err := m.HandleEthUnbond(session, claim)
		if err != nil {
			log.Printf("failed to handle EthUnbonding for topic: %v", utils.EthUnbond)
			return err
		}
	}
	return nil
}

func (m MsgHandler) HandleEthUnbond(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	config := utils.Config()
	producer := utils.NewProducer(m.KafkaConfig.Brokers, config)
	defer func() {
		err := producer.Close()
		if err != nil {
			log.Printf("failed to close producer in topic: %v", utils.EthUnbond)
		}
	}()
	var kafkaMsg *sarama.ConsumerMessage
	defer func() {
		if kafkaMsg != nil {
			session.MarkMessage(kafkaMsg, "")
		}
	}()
	var sum = sdk.NewInt(0)
	for kafkaMsg := range claim.Messages() {
		if kafkaMsg == nil {
			return errors.New("kafka returned nil message")
		}
		var msg sdk.Msg
		err := m.ProtoCodec.UnmarshalInterface(kafkaMsg.Value, &msg)
		if err != nil {
			log.Printf("proto failed to unmarshal")
		}
		switch txMsg := msg.(type) {
		case *bankTypes.MsgSend:
			sum = sum.Add(txMsg.Amount.AmountOf(m.KafkaConfig.Denom))
		default:
			log.Printf("Unexpected type found in topic: %v", utils.EthUnbond)
		}
	}

	unbondMsg := &stakingTypes.MsgUndelegate{
		DelegatorAddress: m.Chain.MustGetAddress().String(),
		ValidatorAddress: constants.Validator1.String(),
		Amount: sdk.Coin{
			Denom:  m.KafkaConfig.Denom,
			Amount: sum,
		},
	}
	msgBytes, err := m.ProtoCodec.MarshalInterface(sdk.Msg(unbondMsg))
	if err != nil {
		return err
	}
	err = utils.ProducerDeliverMessage(msgBytes, utils.ToTendermint, producer)
	if err != nil {
		log.Printf("failed to produce message from topic %v to %v", utils.EthUnbond, utils.ToTendermint)
		return err
	}

	return nil
}

// Handlers of message types
func (m MsgHandler) HandleTopicMsgs(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim, batchSize int,
	handle func([]sarama.ConsumerMessage, *codec.ProtoCodec, *relayer.Chain, *ethclient.Client) error) error {
	msgs := make([]sarama.ConsumerMessage, 0, batchSize)
	for {
		kafkaMsg := <-claim.Messages()
		if kafkaMsg == nil {
			return errors.New("kafka returned nil message")
		}
		log.Printf("Message topic:%q partition:%d offset:%d\n", kafkaMsg.Topic, kafkaMsg.Partition, kafkaMsg.Offset)

		ok, err := BatchAndHandle(&msgs, *kafkaMsg, m.ProtoCodec, m.Chain, m.EthClient, handle)
		if ok && err == nil {
			session.MarkMessage(kafkaMsg, "")
			return nil
		}
		if err != nil {
			return err
		}
	}
}

// BatchAndHandle :
func BatchAndHandle(kafkaMsgs *[]sarama.ConsumerMessage, kafkaMsg sarama.ConsumerMessage,
	protoCodec *codec.ProtoCodec, chain *relayer.Chain, ethClient *ethclient.Client,
	handle func([]sarama.ConsumerMessage, *codec.ProtoCodec, *relayer.Chain, *ethclient.Client) error) (bool, error) {
	*kafkaMsgs = append(*kafkaMsgs, kafkaMsg)
	if len(*kafkaMsgs) == cap(*kafkaMsgs) {
		err := handle(*kafkaMsgs, protoCodec, chain, ethClient)
		if err != nil {
			return false, err
		}
		*kafkaMsgs = (*kafkaMsgs)[:0]
		return true, nil
	}
	return false, nil
}

func ConvertKafkaMsgsToSDKMsg(kafkaMsgs []sarama.ConsumerMessage, protoCodec *codec.ProtoCodec) ([]sdk.Msg, error) {
	var msgs []sdk.Msg
	for _, kafkaMsg := range kafkaMsgs {
		var msg sdk.Msg
		err := protoCodec.UnmarshalInterface(kafkaMsg.Value, &msg)
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, msg)
	}
	return msgs, nil
}

func ConvertKafkaMsgsToEthMsg(kafkaMsgs []sarama.ConsumerMessage) ([]ethereum.EthTxMsg, error) {
	var msgs []ethereum.EthTxMsg
	for _, kafkaMsg := range kafkaMsgs {
		var msg ethereum.EthTxMsg
		err := json.Unmarshal(kafkaMsg.Value, &msg)
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, msg)
	}
	return msgs, nil
}

// SendBatchToEth : Handling of msgSend
func SendBatchToEth(kafkaMsgs []sarama.ConsumerMessage, _ *codec.ProtoCodec, _ *relayer.Chain, ethClient *ethclient.Client) error {
	msgs, err := ConvertKafkaMsgsToEthMsg(kafkaMsgs)
	if err != nil {
		return err
	}
	log.Printf("batched messages to send to ETH: %v", msgs)

	hash, err := ethereum.SendTxToEth(ethClient, msgs, constants.EthGasLimit)
	if err != nil {
		log.Printf("error occuerd in eth transaction: %v", err)
		return err
	}
	log.Printf("sent message to eth with hash: %v ", hash)
	return nil
}

// SendBatchToTendermint :
func SendBatchToTendermint(kafkaMsgs []sarama.ConsumerMessage, protoCodec *codec.ProtoCodec, chain *relayer.Chain, _ *ethclient.Client) error {
	msgs, err := ConvertKafkaMsgsToSDKMsg(kafkaMsgs, protoCodec)
	if err != nil {
		return err
	}
	log.Printf("batched messages to send to Tendermint: %v", msgs)

	response, ok, err := chain.SendMsgs(msgs)
	if err != nil {
		log.Printf("error occured while send to Tendermint:%v: ", err)
		return err
	}
	log.Printf("response: %v, ok: %v", response, ok)
	return nil
}
