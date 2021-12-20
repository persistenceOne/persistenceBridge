/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package handler

import (
	"github.com/Shopify/sarama"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/relayer/relayer"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/persistenceOne/persistenceBridge/utilities/logging"
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
			logging.Error("failed batch and handle for topic ToEth with bridgeErr:", err)
			return err
		}
	case utils.ToTendermint:
		err := m.HandleToTendermint(session, claim)
		if err != nil {
			logging.Error("failed to handle for topic ToTendermint with bridgeErr:", err)
			return err
		}
	case utils.EthUnbond:
		err := m.HandleEthUnbond(session, claim)
		if err != nil {
			logging.Error("failed to handle for topic EthUnbond with bridgeErr:", err)
			return err
		}
	case utils.MsgSend:
		err := m.HandleMsgSend(session, claim)
		if err != nil {
			logging.Error("failed to handle MsgSend for topic MsgSend with bridgeErr:", err)
			return err
		}
	case utils.MsgDelegate:
		err := m.HandleMsgDelegate(session, claim)
		if err != nil {
			logging.Error("failed to handle MsgDelegate for topic MsgDelegate with bridgeErr:", err)
			return err
		}
	case utils.MsgUnbond:
		err := m.HandleMsgUnbond(session, claim)
		if err != nil {
			logging.Error("failed to handle for topic MsgUnbond with bridgeErr:", err)
			return err
		}
	case utils.Redelegate:
		err := m.HandleRelegate(session, claim)
		if err != nil {
			logging.Error("failed to handle for topic Redelegate with bridgeErr:", err)
			return err
		}
	case utils.RetryTendermint:
		err := m.HandleRetryTendermint(session, claim)
		if err != nil {
			logging.Error("failed to handle for topic RetryTendermint with bridgeErr:", err)
			return err
		}
	}

	return nil
}
