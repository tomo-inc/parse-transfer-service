package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type TGBot struct {
	tgBotApi *tgbotapi.BotAPI
	chatId   int64
}

func NewTGBot(
	token string,
	chatId int64,
) *TGBot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal("new bot api failed")
	}
	return &TGBot{
		tgBotApi: bot,
		chatId:   chatId,
	}
}

func (t *TGBot) Send(_ context.Context, msg Msg) (err error) {
	text := fmt.Sprintf("%s <strong> %s </strong> \n\n", t.getLevel(msg.Level), msg.Title)

	for _, v := range msg.Data {
		text = fmt.Sprintf("%s%s : %v \n", text, v.Key, v.Value)
	}

	message := tgbotapi.NewMessage(t.chatId, text)
	message.ParseMode = tgbotapi.ModeHTML
	_, err = t.tgBotApi.Send(message)
	return
}

func (s *TGBot) getLevel(level Level) string {
	var icon string
	switch level {
	case Info:
		icon = "✅"
	case Warning:
		icon = "⚠️"
	case Error:
		icon = "🆘"
	default:
		icon = "🔔"
	}
	return icon
}
