package logging

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"

	"github.com/persistenceOne/persistenceBridge/application/configuration"
)

var bot *tb.Bot
var errorPrefix = []interface{}{"[ERROR]"}
var shutDownWithErrorPrefix = []interface{}{"[SHUTDOWN_WITH_ERROR]"}
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

func Error(err ...interface{}) {
	log.Println(append(errorPrefix, err...)...)
	_ = sendMessage("ERROR:\n" + fmt.Sprintln(err...))
}

func ShutdownWithError(err ...interface{}) {
	log.Println(append(shutDownWithErrorPrefix, err...)...)
	_ = sendMessage("SHUTDOWN_WITH_ERROR:\n" + fmt.Sprintln(err...))
}

func Warn(warn ...interface{}) {
	log.Println(append(warnPrefix, warn...)...)
	_ = sendMessage("WARNING:\n" + fmt.Sprintln(warn...))
}

func Info(info ...interface{}) {
	log.Println(append(infoPrefix, info...)...)
}

func Debug(debug ...interface{}) {
	log.Println(append(debugPrefix, debug...)...)
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
