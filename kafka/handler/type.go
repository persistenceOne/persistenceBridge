package handler

import (
	"github.com/Shopify/sarama"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/relayer/relayer"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"log"
)

type MsgHandler struct {
	ProtoCodec      *codec.ProtoCodec
	Chain           *relayer.Chain
	EthClient       *ethclient.Client
	Count           int
	WithdrawRewards bool
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
			log.Printf("failed to handle MsgSend for topic: %v with error %v\n", utils.MsgSend, err)
			return err
		}
	case utils.MsgDelegate:
		err := m.HandleMsgDelegate(session, claim)
		if err != nil {
			log.Printf("failed to handle MsgDelegate for topic: %v with error %v\n", utils.MsgDelegate, err)
			return err
		}
	case utils.MsgUnbond:
		err := m.HandleMsgUnbond(session, claim)
		if err != nil {
			log.Printf("failed to handle MsgUnbond for topic: %v with error %v\n", utils.MsgUnbond, err)
			return err
		}
	case utils.Redelegate:
		err := m.HandleRelegate(session, claim)
		if err != nil {
			log.Printf("failed to handle Redelegate for topic: %v with error %v\n", utils.Redelegate, err)
			return err
		}
	case utils.RetryTendermint:
		err := m.HandleRetryTendermint(session, claim)
		if err != nil {
			log.Printf("failed to handle for topic: %v with error %v", utils.RetryTendermint, err)
			return err
		}
	}

	return nil
}
