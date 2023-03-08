package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func getResponse(message *tgbotapi.Message) (string, bool) {
	if message == nil {
		return "", false
	}
	if message.Text == "/start" {
		return fmt.Sprintf("Your id is %d", message.Chat.ID), true
	}
	return "Im alive :)", true
}

func getResponseFromMessageData(data *MessageData) string {
	if data.Title == "" {
		return data.Content
	}
	return fmt.Sprintf("*%s*\n\n%s", data.Title, data.Content)
}
