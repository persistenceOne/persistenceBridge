package handler

import (
	"encoding/json"
	"errors"
	"github.com/Shopify/sarama"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/relayer/relayer"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/persistenceOne/persistenceBridge/application"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
	ethereum2 "github.com/persistenceOne/persistenceBridge/ethereum"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/tendermint"
	"log"
)

type MsgHandler struct {
	PstakeConfig configuration.Config
	ProtoCodec   *codec.ProtoCodec
	Chain        *relayer.Chain
	EthClient    *ethclient.Client
	Count        int
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
		err := m.HandleTopicMsgs(session, claim, m.PstakeConfig.Kafka.ToEth.BatchSize, SendBatchToEth)
		if err != nil {
			log.Printf("failed batch and handle for topic: %v with error %v", utils.ToEth, err)
			return err
		}
	case utils.ToTendermint:
		err := m.HandleTopicMsgs(session, claim, m.PstakeConfig.Kafka.ToTendermint.BatchSize, SendBatchToTendermint)
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
	case utils.MsgSend:
		err := m.HandleMsgSend(session, claim)
		if err != nil {
			log.Printf("failed to handle EthUnbonding for topic: %v", utils.MsgSend)
			return err
		}
	case utils.MsgDelegate:
		err := m.HandleMsgDelegate(session, claim)
		if err != nil {
			log.Printf("failed to handle EthUnbonding for topic: %v", utils.MsgDelegate)
			return err
		}
	case utils.MsgUnbond:
		err := m.HandleMsgUnbond(session, claim)
		if err != nil {
			log.Printf("failed to handle EthUnbonding for topic: %v", utils.MsgUnbond)
			return err
		}
	}
	return nil
}

func (m MsgHandler) HandleEthUnbond(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	saramaConfig := utils.SaramaConfig()
	producer := utils.NewProducer(m.PstakeConfig.Kafka.Brokers, saramaConfig)
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
			sum = sum.Add(txMsg.Amount.AmountOf(application.GetAppConfiguration().PStakeDenom))
		default:
			log.Printf("Unexpected type found in topic: %v", utils.EthUnbond)
		}
	}

	if sum != sdk.NewInt(0) {
		// TODO consider multiple validators
		unbondMsg := &stakingTypes.MsgUndelegate{
			DelegatorAddress: m.Chain.MustGetAddress().String(),
			ValidatorAddress: constants2.Validator1.String(),
			Amount: sdk.Coin{
				Denom:  application.GetAppConfiguration().PStakeDenom,
				Amount: sum,
			},
		}
		msgBytes, err := m.ProtoCodec.MarshalInterface(sdk.Msg(unbondMsg))
		if err != nil {
			return err
		}
		err = utils.ProducerDeliverMessage(msgBytes, utils.MsgUnbond, producer)
		if err != nil {
			log.Printf("failed to produce message from topic %v to %v", utils.EthUnbond, utils.ToTendermint)
			return err
		}
	}

	return nil
}

// HandleTopicMsgs Handlers of message types
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

func ConvertKafkaMsgsToEthMsg(kafkaMsgs []sarama.ConsumerMessage) ([]ethereum2.EthTxMsg, error) {
	var msgs []ethereum2.EthTxMsg
	for _, kafkaMsg := range kafkaMsgs {
		var msg ethereum2.EthTxMsg
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

	hash, err := ethereum2.SendTxToEth(ethClient, msgs, application.GetAppConfiguration().EthGasLimit)
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

func (m MsgHandler) HandleMsgSend(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	config := utils.SaramaConfig()
	producer := utils.NewProducer(m.PstakeConfig.Kafka.Brokers, config)
	defer func() {
		err := producer.Close()
		if err != nil {
			log.Printf("failed to close producer in topic: %v", utils.MsgSend)
		}
	}()

	messagesLength := len(claim.Messages())
	loop := messagesLength
	if messagesLength > m.PstakeConfig.Kafka.ToTendermint.BatchSize-m.Count {
		loop = m.PstakeConfig.Kafka.ToTendermint.BatchSize - m.Count
	}
	if messagesLength > 0 {
		var msgs [][]byte
		// TODO add msg withdraw rewards from multiple validators.
		if tendermint.AddressIsDelegatorToValidator(m.Chain.MustGetAddress().String(), constants2.Validator1.String(), m.Chain) {
			withdrawRewardsMsg := &distributionTypes.MsgWithdrawDelegatorReward{
				DelegatorAddress: m.Chain.MustGetAddress().String(),
				ValidatorAddress: constants2.Validator1.String(),
			}
			withdrawRewardsMsgBytes, err := m.ProtoCodec.MarshalInterface(sdk.Msg(withdrawRewardsMsg))
			if err != nil {
				log.Printf("Failed to Marshal WithdrawMessage: Error: %v", err)
			} else {
				msgs = append(msgs, withdrawRewardsMsgBytes)
				loop = loop - 1
			}
		}

		var kafkaMsg *sarama.ConsumerMessage
		for i := 0; i < loop; i++ {
			kafkaMsg = <-claim.Messages()
			if kafkaMsg == nil {
				return errors.New("kafka returned nil message")
			}
			msgs = append(msgs, kafkaMsg.Value)
		}
		if len(msgs) > 0 {
			err := utils.ProducerDeliverMessages(msgs, utils.ToTendermint, producer)
			if err != nil {
				log.Printf("error in handler for topic %v, failed to produce to queue", utils.MsgSend)
			}
			session.MarkMessage(kafkaMsg, "")
			return err
		}
	}
	m.Count += loop
	return nil
}
func (m MsgHandler) HandleMsgDelegate(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	config := utils.SaramaConfig()
	producer := utils.NewProducer(m.PstakeConfig.Kafka.Brokers, config)
	defer func() {
		err := producer.Close()
		if err != nil {
			log.Printf("failed to close producer in topic: %v", utils.MsgDelegate)
		}
	}()
	messagesLength := len(claim.Messages())
	if messagesLength > 0 {
		var msgs [][]byte
		var kafkaMsg *sarama.ConsumerMessage
		for i := 0; i < messagesLength; i++ {
			kafkaMsg = <-claim.Messages()
			if kafkaMsg == nil {
				return errors.New("kafka returned nil message")
			}

			msgs = append(msgs, kafkaMsg.Value)
		}
		if len(msgs) > 0 {
			err := utils.ProducerDeliverMessages(msgs, utils.ToTendermint, producer)
			if err != nil {
				log.Printf("error in handler for topic %v, failed to produce to queue", utils.MsgDelegate)
			}
			session.MarkMessage(kafkaMsg, "")
			return err
		}
	}
	m.Count += messagesLength
	return nil
}
func (m MsgHandler) HandleMsgUnbond(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	config := utils.SaramaConfig()
	producer := utils.NewProducer(m.PstakeConfig.Kafka.Brokers, config)
	defer func() {
		err := producer.Close()
		if err != nil {
			log.Printf("failed to close producer in topic: %v", utils.MsgUnbond)
		}
	}()
	messagesLength := len(claim.Messages())
	if messagesLength > 0 {
		var msgs [][]byte
		var kafkaMsg *sarama.ConsumerMessage
		for i := 0; i < messagesLength; i++ {
			kafkaMsg = <-claim.Messages()
			if kafkaMsg == nil {
				return errors.New("kafka returned nil message")
			}
			msgs = append(msgs, kafkaMsg.Value)
		}
		if len(msgs) > 0 {
			err := utils.ProducerDeliverMessages(msgs, utils.ToTendermint, producer)
			if err != nil {
				log.Printf("error in handler for topic %v, failed to produce to queue", utils.MsgUnbond)
			}
			session.MarkMessage(kafkaMsg, "")
			return err
		}
	}
	m.Count += messagesLength
	return nil
}
