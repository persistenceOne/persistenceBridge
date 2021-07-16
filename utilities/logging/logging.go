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
			logrus.Fatalln("ERROR while setting up bot: %s\n", err.Error())
		}
	}
	if configuration.GetAppConfig().TelegramBot.ChatID == 0 {
		err = fmt.Errorf("invalid bot's chat id")
	} else {
		err = sendMessage("pBridge bot initialized")
	}
	return err
}

func Error(err ...interface{}) {
	logrus.Errorln(err...)
	message := fmt.Sprintln(err...)
	_ = sendMessage("ERROR:\n" + message)
}

func Warn(warn ...interface{}) {
	logrus.Warnln(warn...)
	message := fmt.Sprintln(warn...)
	_ = sendMessage("WARNING:\n" + message)
}

func Info(info ...interface{}) {
	logrus.Infoln(info...)
}

func Debug(debug ...interface{}) {
	logrus.Debugln(debug...)
}

func Fatal(err ...interface{}) {
	message := fmt.Sprintln(err...)
	_ = sendMessage("FATAL:\n" + message)
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
