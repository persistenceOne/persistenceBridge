package logging

import (
	"fmt"
	"github.com/persistenceOne/persistenceBridge/application/configuration"
	"github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
)

var bot *tb.Bot

func InitializeBotAndLog(logLevel string) (err error) {
	ll, err := logrus.ParseLevel(logLevel)
	if err != nil {
		panic(err)
	}
	logrus.SetLevel(ll)
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)
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
	logrus.Errorln(err...)
	_ = sendMessage("ERROR:\n" + fmt.Sprintln(err...))
}

func Warn(warn ...interface{}) {
	logrus.Warnln(warn...)
	_ = sendMessage("WARNING:\n" + fmt.Sprintln(warn...))
}

func Info(info ...interface{}) {
	logrus.Infoln(info...)
}

func Debug(debug ...interface{}) {
	logrus.Debugln(debug...)
}

func Fatal(err ...interface{}) {
	_ = sendMessage("FATAL:\n" + fmt.Sprintln(err...))
	logrus.Fatalln(err...)
}

func sendMessage(message string) error {
	if bot != nil {
		_, err := bot.Send(tb.ChatID(configuration.GetAppConfig().TelegramBot.ChatID), message)
		if err != nil {
			logrus.Errorln("Bot send message error: %s\n", err.Error())
			return err
		}
	}
	return nil
}
