package handler

import (
	"github.com/Shopify/sarama"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/relayer/relayer"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
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
		err := m.HandleToEth(session, claim)
		if err != nil {
			log.Printf("failed batch and handle for topic: %v with error %v\n", utils.ToEth, err)
			return err
		}
	case utils.ToTendermint:
		err := m.HandleToTendermint(session, claim)
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
