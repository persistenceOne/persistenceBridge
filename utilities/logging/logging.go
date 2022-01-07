/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package logging

import (
	"fmt"
	"log"

	tb "gopkg.in/tucnak/telebot.v2"

	"github.com/persistenceOne/persistenceBridge/application/configuration"
)

// nolint fixme: move into a config or a proper structure type
// nolint: gochecknoglobals
var (
	bot       *tb.Bot
	showDebug bool
)

const (
	errorPrefix = "[ERROR]"
	warnPrefix  = "[WARNING]"
	infoPrefix  = "[INFO]"
	debugPrefix = "[DEBUG]"
	fatalPrefix = "[FATAL]"
)

func InitializeBot() (err error) {
	if configuration.GetAppConfig().TelegramBot.Token != "" {
		bot, err = tb.NewBot(tb.Settings{Token: configuration.GetAppConfig().TelegramBot.Token})
		if err != nil {
			return err
		}

		if configuration.GetAppConfig().TelegramBot.ChatID != 0 {
			Info("Sending initializing bot message...")

			err = sendMessage("pBridge bot initialized")
			if err != nil {
				return err
			}
		}
	}

	return err
}

// fixme introduce a proper logger and put it into a context
func ShowDebugLog(d bool) {
	showDebug = d
}

func Error(errs ...interface{}) {
	logPrefix(errorPrefix, errs)
	_ = sendMessage("ERROR:\n" + fmt.Sprintln(errs...))
}

func Errorf(format string, errs ...interface{}) {
	logfPrefix(format, errorPrefix, errs)
	_ = sendMessage("ERROR:\n" + fmt.Sprintln(errs...))
}

func Warn(warns ...interface{}) {
	logPrefix(warnPrefix, warns)
	_ = sendMessage("WARNING:\n" + fmt.Sprintln(warns...))
}

func Warnf(format string, errs ...interface{}) {
	logfPrefix(format, warnPrefix, errs)
	_ = sendMessage("WARNING:\n" + fmt.Sprintln(errs...))
}

func Info(infos ...interface{}) {
	logPrefix(infoPrefix, infos)
}

func Infof(format string, infos ...interface{}) {
	logfPrefix(format, infoPrefix, infos)
}

func Debug(debugs ...interface{}) {
	if showDebug {
		logPrefix(debugPrefix, debugs)
	}
}

func Debugf(format string, debugs ...interface{}) {
	if showDebug {
		logfPrefix(format, debugPrefix, debugs)
	}
}

func Fatal(errs ...interface{}) {
	_ = sendMessage("FATAL:\n" + fmt.Sprintln(errs...))
	logPrefix(fatalPrefix, errs)
}

func Fatalf(format string, errs ...interface{}) {
	_ = sendMessage("FATAL:\n" + fmt.Sprintln(errs...))
	logfPrefix(format, fatalPrefix, errs)
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

func logPrefix(prefix string, errs ...interface{}) {
	log.Println(append([]interface{}{prefix}, errs...)...)
}

func logfPrefix(format, prefix string, errs ...interface{}) {
	log.Println(prefix, fmt.Sprintf(format, errs...))
}
