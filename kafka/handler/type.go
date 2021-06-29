package handler

import (
	"encoding/json"
	"errors"
	"github.com/Shopify/sarama"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/relayer/relayer"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	constants2 "github.com/persistenceOne/persistenceBridge/application/constants"
	ethereum2 "github.com/persistenceOne/persistenceBridge/ethereum"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/tendermint"
	"log"
	"time"
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
		err := m.HandleToEth(session, claim, m.PstakeConfig.Kafka.ToEth.BatchSize)
		if err != nil {
			log.Printf("failed batch and handle for topic: %v with error %v\n", utils.ToEth, err)
			return err
		}
	case utils.ToTendermint:
		err := m.HandleToTendermint(session, claim, m.PstakeConfig.Kafka.ToTendermint.BatchSize)
		if err != nil {
			log.Printf("failed to handle for topic: %v with error %v\n", utils.ToTendermint, err)
			return err
		}
	case utils.EthUnbond:
		err := m.HandleEthUnbond(session, claim)
		if err != nil {
			log.Printf("failed to handle EthUnbonding for topic: %v\n", utils.EthUnbond)
			return err
		}
	case utils.MsgSend:
		err := m.HandleMsgSend(session, claim)
		if err != nil {
			log.Printf("failed to handle EthUnbonding for topic: %v\n", utils.MsgSend)
			return err
		}
	case utils.MsgDelegate:
		err := m.HandleMsgDelegate(session, claim)
		if err != nil {
			log.Printf("failed to handle EthUnbonding for topic: %v\n", utils.MsgDelegate)
			return err
		}
	case utils.MsgUnbond:
		err := m.HandleMsgUnbond(session, claim)
		if err != nil {
			log.Printf("failed to handle EthUnbonding for topic: %v\n", utils.MsgUnbond)
			return err
		}
	}
	return nil
}

func (m MsgHandler) HandleToEth(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim, batchSize int) error {
	msgs := make([]sarama.ConsumerMessage, 0, batchSize)
	for {
		kafkaMsg := <-claim.Messages()
		if kafkaMsg == nil {
			return errors.New("kafka returned nil message")
		}
		log.Printf("Message topic:%q partition:%d offset:%d\n", kafkaMsg.Topic, kafkaMsg.Partition, kafkaMsg.Offset)

		ok, err := BatchAndHandleEthereum(&msgs, *kafkaMsg, m)
		if ok && err == nil {
			session.MarkMessage(kafkaMsg, "")
			return nil
		}
		if err != nil {
			return err
		}
	}
}
func (m MsgHandler) HandleToTendermint(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim, batchSize int) error {
	var kafkaMsgs []sarama.ConsumerMessage
	claimMsgChan := claim.Messages()
	closeChan := make(chan bool, 1)
	ticker := time.Tick(1 * time.Second)
	var kafkaMsg *sarama.ConsumerMessage
	var ok bool
	for {
		select {
		case <-closeChan:
			if len(kafkaMsgs) == 0 {
				return nil
			}
			if kafkaMsg == nil {
				return errors.New("kafka returned nil message")
			}
			err := SendBatchToTendermint(kafkaMsgs, m)
			if err != nil {
				return err
			}
			session.MarkMessage(kafkaMsg, "")
			return nil
		case <-ticker:
			if len(kafkaMsgs) != 0 {
				AddToBufferedChannelIfCapacityPermits(closeChan, true)
			}
		case kafkaMsg, ok = <-claimMsgChan:
			if ok {
				kafkaMsgs = append(kafkaMsgs, *kafkaMsg)
				if len(kafkaMsgs) == batchSize {
					AddToBufferedChannelIfCapacityPermits(closeChan, true)
				} else if len(kafkaMsgs) > batchSize {
					log.Printf("Select tried to batch more messages in handler: %v ,not "+
						"comitting offset, %v", utils.ToTendermint, kafkaMsg.Offset)
					return nil
				}
			} else {
				AddToBufferedChannelIfCapacityPermits(closeChan, true)
			}
		}
	}
}

func (m MsgHandler) HandleEthUnbond(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	saramaConfig := utils.SaramaConfig()
	producer := utils.NewProducer(m.PstakeConfig.Kafka.Brokers, saramaConfig)
	defer func() {
		err := producer.Close()
		if err != nil {
			log.Printf("failed to close producer in topic: %v\n", utils.EthUnbond)
		}
	}()
	var kafkaMsg *sarama.ConsumerMessage
	defer func() {
		if kafkaMsg != nil {
			session.MarkMessage(kafkaMsg, "")
		}
	}()
	var sum = sdk.NewInt(0)
	messagesLength := len(claim.Messages())

	for i := 0; i < messagesLength; i++ {
		kafkaMsg = <-claim.Messages()
		if kafkaMsg == nil {
			return errors.New("kafka returned nil message")
		}

		var msg sdk.Msg
		err := m.ProtoCodec.UnmarshalInterface(kafkaMsg.Value, &msg)
		if err != nil {
			log.Printf("proto failed to unmarshal\n")
		}
		switch txMsg := msg.(type) {
		case *stakingTypes.MsgUndelegate:
			sum = sum.Add(txMsg.Amount.Amount)
		default:
			log.Printf("Unexpected type found in topic: %v\n", utils.EthUnbond)
		}

	}

	if sum.GT(sdk.NewInt(0)) {
		// TODO consider multiple validators
		unbondMsg := &stakingTypes.MsgUndelegate{
			DelegatorAddress: m.Chain.MustGetAddress().String(),
			ValidatorAddress: constants2.Validator1.String(),
			Amount: sdk.Coin{
				Denom:  configuration.GetAppConfiguration().PStakeDenom,
				Amount: sum,
			},
		}
		msgBytes, err := m.ProtoCodec.MarshalInterface(sdk.Msg(unbondMsg))
		if err != nil {
			return err
		}

		err = utils.ProducerDeliverMessage(msgBytes, utils.MsgUnbond, producer)
		if err != nil {
			log.Printf("failed to produce message from topic %v to %v\n", utils.EthUnbond, utils.MsgUnbond)
			return err
		}
	}
	return nil
}

func (m MsgHandler) HandleMsgSend(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	config := utils.SaramaConfig()
	producer := utils.NewProducer(m.PstakeConfig.Kafka.Brokers, config)
	defer func() {
		err := producer.Close()
		if err != nil {
			log.Printf("failed to close producer in topic: %v\n", utils.MsgSend)
		}
	}()

	messagesLength := len(claim.Messages())
	loop := messagesLength
	if m.PstakeConfig.Kafka.ToTendermint.BatchSize-m.Count <= 0 {
		return nil
	}
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
				log.Printf("Failed to Marshal WithdrawMessage: Error: %v\n", err)
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
			session.MarkMessage(kafkaMsg, "")
			if err != nil {
				log.Printf("error in handler for topic %v, failed to produce to queue\n", utils.MsgSend)
				return err
			}
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
			log.Printf("failed to close producer in topic: %v\n", utils.MsgDelegate)
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
			session.MarkMessage(kafkaMsg, "")
			if err != nil {
				log.Printf("error in handler for topic %v, failed to produce to queue\n", utils.MsgDelegate)
				return err
			}
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
			log.Printf("failed to close producer in topic: %v\n", utils.MsgUnbond)
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
			session.MarkMessage(kafkaMsg, "")
			if err != nil {
				log.Printf("error in handler for topic %v, failed to produce to queue\n", utils.MsgUnbond)
				return err
			}
		}
	}
	m.Count += messagesLength
	return nil
}

//// BatchAndHandle :
//func BatchAndHandleTendermint(kafkaMsgs *[]sarama.ConsumerMessage, kafkaMsg sarama.ConsumerMessage, msgHandler MsgHandler) (bool, error) {
//	*kafkaMsgs = append(*kafkaMsgs, kafkaMsg)
//	if len(*kafkaMsgs) == cap(*kafkaMsgs) {
//		err := SendBatchToTendermint(*kafkaMsgs, msgHandler)
//		if err != nil {
//			return false, err
//		}
//		*kafkaMsgs = (*kafkaMsgs)[:0]
//		return true, nil
//	}
//	return false, nil
//}

// BatchAndHandle :
func BatchAndHandleEthereum(kafkaMsgs *[]sarama.ConsumerMessage, kafkaMsg sarama.ConsumerMessage, msgHandler MsgHandler) (bool, error) {
	*kafkaMsgs = append(*kafkaMsgs, kafkaMsg)
	if len(*kafkaMsgs) == cap(*kafkaMsgs) {
		err := SendBatchToEth(*kafkaMsgs, msgHandler)
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
func SendBatchToEth(kafkaMsgs []sarama.ConsumerMessage, handler MsgHandler) error {
	msgs, err := ConvertKafkaMsgsToEthMsg(kafkaMsgs)
	if err != nil {
		return err
	}
	log.Printf("batched messages to send to ETH: %v\n", msgs)

	hash, err := ethereum2.SendTxToEth(handler.EthClient, msgs, configuration.GetAppConfiguration().Ethereum.EthGasLimit)
	if err != nil {
		log.Printf("error occuerd in sending eth transaction: %v\n", err)
		return err
	}
	log.Printf("Broadcasted Eth Tx hash: %s\n", hash)
	return nil
}

// SendBatchToTendermint :
func SendBatchToTendermint(kafkaMsgs []sarama.ConsumerMessage, handler MsgHandler) error {
	msgs, err := ConvertKafkaMsgsToSDKMsg(kafkaMsgs, handler.ProtoCodec)
	if err != nil {
		return err
	}
	log.Printf("batched messages to send to Tendermint: %v\n", msgs)

	response, ok, err := handler.Chain.SendMsgs(msgs)
	if err != nil {
		log.Printf("error occured while send to Tendermint:%v\n", err)
		return err
	}
	if !ok {
		config := utils.SaramaConfig()
		producer := utils.NewProducer(handler.PstakeConfig.Kafka.Brokers, config)
		defer func() {
			err := producer.Close()
			if err != nil {
				log.Printf("failed to close producer in topic: %v\n", utils.MsgUnbond)
			}
		}()

		for _, msg := range msgs {
			msgBytes, err := handler.ProtoCodec.MarshalInterface(msg)
			if err != nil {
				log.Printf("Retry txs: Failed to Marshal ToTendermint Retry msg: Error: %v\n", err)
			}
			err = utils.ProducerDeliverMessage(msgBytes, utils.ToTendermint, producer)
			if err != nil {
				log.Printf("Retry txs: Failed to add msg to kafka queue: %s\n", err.Error())
			}
			log.Printf("Retry txs: Produced to kafka: %v, for topic %v\n", msg.String(), utils.ToTendermint)
		}
	}
	log.Printf("Broadcasted Tendermint TX HASH: %s, ok: %v\n", response.TxHash, ok)
	return nil
}

func AddToBufferedChannelIfCapacityPermits(channel chan bool, data bool) {
	if len(channel) < cap(channel) {
		channel <- data
	}
}
