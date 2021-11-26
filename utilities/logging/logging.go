package logging

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/persistenceOne/persistenceBridge/kafka/utils"
	"github.com/pkg/errors"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"math/big"

	"github.com/persistenceOne/persistenceBridge/application/configuration"
)

var bot *tb.Bot
var showDebug bool
var errorPrefix = []interface{}{"[ERROR]"}
var warnPrefix = []interface{}{"[WARNING]"}
var infoPrefix = []interface{}{"[INFO]"}
var debugPrefix = []interface{}{"[DEBUG]"}
var fatalPrefix = []interface{}{"[FATAL]"}

func InitializeBot() (err error) {
	if configuration.GetAppConfig().TelegramBot.Token != "" {
		bot, err = tb.NewBot(tb.Settings{Token: configuration.GetAppConfig().TelegramBot.Token})
		if err != nil {
			return err
		}
		if configuration.GetAppConfig().TelegramBot.ChatID != 0 {
			Info("Sending initialising bot message...")
			err = sendMessage("pBridge bot initialized")
			if err != nil {
				return err
			}
		}
	}
	return err
}

func ShowDebugLog(d bool) {
	showDebug = d
}

func Error(err ...interface{}) {
	log.Println(append(errorPrefix, err...)...)
	_ = sendMessage("ERROR:\n" + fmt.Sprintln(err...))
}

func Warn(warn ...interface{}) {
	log.Println(append(warnPrefix, warn...)...)
	_ = sendMessage("WARNING:\n" + fmt.Sprintln(warn...))
}

func Info(info ...interface{}) {
	log.Println(append(infoPrefix, info...)...)
}

func Debug(debug ...interface{}) {
	if showDebug {
		log.Println(append(debugPrefix, debug...)...)
	}
}

func Fatal(err ...interface{}) {
	_ = sendMessage("FATAL:\n" + fmt.Sprintln(err...))
	log.Fatalln(append(fatalPrefix, err...)...)
}

func sendMessage(message string) error {
	if bot != nil {
		_, err := bot.Send(tb.ChatID(configuration.GetAppConfig().TelegramBot.ChatID), message)
		if err != nil {
			log.Println("Bot send message error:", err.Error())
			return err
		}
	}
	return nil
}

func ErrorReason(ctx context.Context, from common.Address, tx *types.Transaction, blockNum *big.Int, ethClient ethclient.Client) error {
	if tx.To().Hash() == common.HexToHash(utils.BlackHole){
		return errors.New( "ZeroAddress not allowed")
	}
	msg := ethereum.CallMsg{
		From:     from,
		To:       tx.To(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Value:    tx.Value(),
		Data:     tx.Data(),
	}

	res, err := ethClient.CallContract(ctx, msg, blockNum)
	if err != nil {
		return errors.Wrap(err, "Cannot get revert reason fatal")
	}
	if len(res) == 0 {
		return errors.New("Out of gas")
	}

	return unpackError(res)
}

var (
	errorSig            = []byte{0x08, 0xc3, 0x79, 0xa0} // Keccak256("Error(string)")[:4]
	abiString, _        = abi.NewType("string", "", nil)
)

func unpackError(result []byte) error {
	if !bytes.Equal(result[:4], errorSig) {
		return errors.New("TX result not of type Error(string)")
	}
	vs, err := abi.Arguments{{Type: abiString}}.UnpackValues(result[4:])
	if err != nil {
		return errors.Wrap(err, "Unpacking revert reason")
	}
	return errors.New(vs[0].(string))
}
