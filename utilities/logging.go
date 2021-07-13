package utilities

import (
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
	"os"
)

var Token string
var ChatId int64

func InitLog() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	log.Infof(Token)
	//Bot = tb.Bot{Token: token}
}

func LogError(err string) {
	log.Errorf(err)
	if Token != "" {
		if ChatId != 0 {
			Bot, _ := tb.NewBot(tb.Settings{Token: Token})
			var ch_id tb.Recipient = tb.ChatID(ChatId)
			_, _ = Bot.Send(ch_id, "Error occurred in pbridge:\n"+err)
		}else {
			log.Fatalf("Please enter the chat id as well")
		}
	}
}
func LogWarning(err string) {
	log.Warnf(err)
}
func LogInfo(err string) {
	log.Infof(err)
}
func LogDebug(err string) {
	log.Debugf(err)
}
func LogFatal(err string) {
	log.Fatalf(err)
	if Token != "" {
		if ChatId != 0 {

			Bot, _ := tb.NewBot(tb.Settings{Token: Token})
			var ch_id tb.Recipient = tb.ChatID(ChatId)
			_, _ = Bot.Send(ch_id, "Error occurred in pbridge:\n"+err)
		}else {
			log.Fatalf("Please enter the chat id as well")
		}
	}
}
