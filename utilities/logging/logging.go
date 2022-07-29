/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceBridge contributors
 SPDX-License-Identifier: Apache-2.0
*/

package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/persistenceOne/persistenceBridge/application/constants"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"net/http"

)

var bot *tb.Bot
var showDebug bool
var errorPrefix = []interface{}{"[ERROR]"}
var warnPrefix = []interface{}{"[WARNING]"}
var infoPrefix = []interface{}{"[INFO]"}
var debugPrefix = []interface{}{"[DEBUG]"}
var fatalPrefix = []interface{}{"[FATAL]"}

func InitializeBot() (err error) {
	if configuration.GetAppConfig().InitSlackBot {
		err = sendSlackMessage("pBridge bot initialized")
		if err != nil {
			return err
		}
	}

	values := map[string]string{"name": "John Doe", "occupation": "gardener"}
	json_data, err := json.Marshal(values)
	_, err = http.Post("https://hooks.slack.com/services" + constants.SlackToken, "application/json",
		bytes.NewBuffer(json_data))

	if err != nil {
		log.Fatal(err)
	}

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

	if err != nil {
		log.Fatal(err)
	}

	return err
}

func ShowDebugLog(d bool) {
	showDebug = d
}

func Error(err ...interface{}) {
	log.Println(append(errorPrefix, err...)...)
	_ = sendMessage("ERROR:\n" + fmt.Sprintln(err...))
	_ = sendSlackMessage("ERROR:\n" + fmt.Sprintln(err...))
}

func Warn(warn ...interface{}) {
	log.Println(append(warnPrefix, warn...)...)
	_ = sendMessage("WARNING:\n" + fmt.Sprintln(warn...))
	_ = sendSlackMessage("WARNING:\n" + fmt.Sprintln(warn...))
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
	_ = sendSlackMessage("FATAL:\n" + fmt.Sprintln(err...))

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

func sendSlackMessage(message string) error {
	values := map[string]string{"text": message}

	jsonData, err := json.Marshal(values)
	_, err = http.Post("https://hooks.slack.com/services"+constants.SlackToken, "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
