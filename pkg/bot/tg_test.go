package bot

import (
	"context"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
)

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("1.com", "http://1.com"),
		tgbotapi.NewInlineKeyboardButtonData("2", "2"),
		tgbotapi.NewInlineKeyboardButtonData("3", "3"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("4", "4"),
		tgbotapi.NewInlineKeyboardButtonData("5", "5"),
		tgbotapi.NewInlineKeyboardButtonData("6", "6"),
	),
)

// see chat group
// curl https://api.telegram.org/bot7792674504:AAEIRU8D_0E2oqUCwYoWGlSILoPhCd6vFjU/getUpdates

func Test_Tg(t *testing.T) {
	bot, err := tgbotapi.NewBotAPI("7792674504:AAEIRU8D_0E2oqUCwYoWGlSILoPhCd6vFjU")
	assert.NoError(t, err)
	msg := tgbotapi.NewMessage(-1002374916275, `
<b>bold</b>, <strong>bold</strong>
<i>italic</i>, <em>italic</em>
<u>underline</u>, <ins>underline</ins>
<s>strikethrough</s>, <strike>strikethrough</strike>, <del>strikethrough</del>
<span class="tg-spoiler">spoiler</span>, <tg-spoiler>spoiler</tg-spoiler>
<b>bold <i>italic bold <s>italic bold strikethrough <span class="tg-spoiler">italic bold strikethrough spoiler</span></s> <u>underline italic bold</u></i> bold</b>
<a href="http://www.example.com/">inline URL</a>
<a href="tg://user?id=123456789">inline mention of a user</a>
<tg-emoji emoji-id="5368324170671202286">👍</tg-emoji>
`)
	msg.ParseMode = tgbotapi.ModeHTML
	//msg.ReplyMarkup = numericKeyboard
	message, err := bot.Send(msg)
	assert.NoError(t, err)
	t.Log(message)
}

func Test_TgBot(t *testing.T) {
	tgbot := NewTGBot("7792674504:AAEIRU8D_0E2oqUCwYoWGlSILoPhCd6vFjU", -1002374916275)
	data := Msg{
		Title: "hahahah",
		Level: Info,
		Data: []*KeyValue{
			{
				Key:   "name",
				Value: "dsadada",
			},
			{
				Key:   "age",
				Value: 12,
			},
			{
				Key:   "school",
				Value: "aaa",
			},
		},
	}
	err := tgbot.Send(context.Background(), data)
	assert.NoError(t, err)
}
