package main

import (
	"time"

	"github.com/ofek1weiss/notifying-bot/bot"
	"github.com/sirupsen/logrus"
)

func panicOnErr(err error) {
	if err != nil {
		logrus.Panic(err)
	}
}

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	config, err := bot.LoadFile[bot.Config]("config.json")
	panicOnErr(err)
	messageDataLoader := bot.MessageDataLoader{
		Path:           config.MessagesPath,
		SampleInterval: 5 * time.Second,
	}
	tgbot := bot.NewBot(config.AccessToken, messageDataLoader.Listen(), config.DefaultChatId)
	panicOnErr(tgbot.Start())
}
