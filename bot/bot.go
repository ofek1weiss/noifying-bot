package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type Bot struct {
	accessToken   string
	messagesIn    <-chan *MessageData
	defaultChatId int64
	botApi        *tgbotapi.BotAPI
}

func NewBot(accessToken string, messagesIn <-chan *MessageData, defaultChatId int64) *Bot {
	return &Bot{
		accessToken:   accessToken,
		messagesIn:    messagesIn,
		defaultChatId: defaultChatId,
	}
}

func (b *Bot) Start() error {
	botApi, err := tgbotapi.NewBotAPI(b.accessToken)
	if err != nil {
		return err
	}
	logrus.Info("Bot connected")
	b.botApi = botApi
	return b.listen()
}

func (b *Bot) listen() error {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updatesIn := b.botApi.GetUpdatesChan(updateConfig)
	for {
		var err error = nil
		select {
		case messageData := <-b.messagesIn:
			logrus.Info("Got message data")
			err = b.handleMessageData(messageData)
		case update := <-updatesIn:
			logrus.Info("Got update from server")
			err = b.handleUpdate(update)
		}
		if err != nil {
			logrus.Error(err)
		}
	}
}

func (b *Bot) handleMessageData(data *MessageData) error {
	if data.ChatID == 0 {
		data.ChatID = b.defaultChatId
	}
	text := getResponseFromMessageData(data)
	message := tgbotapi.NewMessage(data.ChatID, text)
	return b.sendMessage(&message)
}

func (b *Bot) handleUpdate(update tgbotapi.Update) error {
	response, ok := getResponse(update.Message)
	if !ok {
		return nil
	}
	message := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	return b.sendMessage(&message)
}

func (b *Bot) sendMessage(message *tgbotapi.MessageConfig) error {
	logrus.Info("Sending message", message.Text)
	message.ParseMode = "markdown"
	_, err := b.botApi.Send(message)
	return err
}
