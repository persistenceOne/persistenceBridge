package utilities

import (
	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
)

var Token string
var ChatId int64
var bot *tb.Bot

func InitLog() {
	//file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.SetOutput(file)
	//log.Infof(Token)

	bot, _ = tb.NewBot(tb.Settings{Token: Token})

}

func LogError(err string) {
	log.Errorf(err)
	sendMessage("Error occurred in pbridge:\n"+err)
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
	sendMessage("Exiting bridge due to a falat error.:\n"+err)
	log.Fatalf(err)
}

func sendMessage(mess string){
	if bot != nil {
		if Token != "" {
			if ChatId != 0 {
				var ch_id tb.Recipient = tb.ChatID(ChatId)
				_, _ = bot.Send(ch_id, mess)
			} else {
				log.Warnf("Chat id not mentioned. Can't alert to telegram.")
			}
		}else {
			log.Warnf("Telegram token not given. Alerts won't work")
		}
	}else {
		log.Warnf("Bot not initiated.")
	}
}
