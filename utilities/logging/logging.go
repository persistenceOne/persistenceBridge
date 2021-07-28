package logging

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"

	"github.com/persistenceOne/persistenceBridge/application/configuration"
)

var bot *tb.Bot

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

func Error(err ...interface{}) {
	log.Println(append([]interface{}{"[ERROR]"}, err...)...)
	_ = sendMessage("ERROR:\n" + fmt.Sprintln(err...))
}

func Warn(warn ...interface{}) {
	log.Println(append([]interface{}{"[WARN]"}, warn...)...)
	_ = sendMessage("WARNING:\n" + fmt.Sprintln(warn...))
}

func Info(info ...interface{}) {
	log.Println(append([]interface{}{"[INFO]"}, info...)...)
}

func Debug(debug ...interface{}) {
	log.Println(append([]interface{}{"[DEBUG]"}, debug...)...)
}

func Fatal(err ...interface{}) {
	_ = sendMessage("FATAL:\n" + fmt.Sprintln(err...))
	log.Println(append([]interface{}{"[FATAL]"}, err...)...)
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
