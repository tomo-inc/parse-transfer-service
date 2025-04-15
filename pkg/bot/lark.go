package bot

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

var larkV2Url = "https://open.larksuite.com/open-apis/bot/v2/hook/"

type LarkBot struct {
	botID  string
	secret string
}

func NewLarkBot(
	botID string,
	secret string,
) *LarkBot {
	return &LarkBot{
		botID:  botID,
		secret: secret,
	}
}

func (t *LarkBot) Send(_ context.Context, msg Msg) (err error) {
	whmsg, err := t.generateMessage(&msg)
	if err != nil {
		return err
	}
	data, err := json.Marshal(whmsg)
	if err != nil {
		return err
	}
	resp, err := http.Post(fmt.Sprintf("%s%s", larkV2Url, t.botID), "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return
}

func (t *LarkBot) generateMessage(msg *Msg) (*WebHookMessage, error) {
	texts := make([]*Text, 0, len(msg.Data))
	elements := make([]*Element, 0, len(texts))
	for _, v := range msg.Data {
		elements = append(elements, &Element{
			Tag: "div",
			Text: &Text{
				Tag:     "lark_md",
				Content: fmt.Sprintf("**%s** : %v", v.Key, v.Value),
			},
		})
	}
	card := &WebhookMessageCard{
		Elements: elements,
		Header: &Header{Title: Text{
			Content: fmt.Sprintf(" %s %s", t.getLevel(msg.Level), msg.Title),
			Tag:     "plain_text",
		}},
	}

	now := time.Now().Unix()
	signature, err := t.sign(now)
	if err != nil {
		return nil, err
	}
	return &WebHookMessage{
		Timestamp: strconv.Itoa(int(now)),
		Sign:      signature,
		MsgType:   "interactive",
		Card:      card,
	}, nil
}

func (t *LarkBot) sign(timestamp int64) (string, error) {
	stringToSign := fmt.Sprintf("%v", timestamp) + "\n" + t.secret

	var data []byte
	h := hmac.New(sha256.New, []byte(stringToSign))
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}

	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature, nil
}

func (s *LarkBot) getLevel(level Level) string {
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

type Text struct {
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

type Element struct {
	Tag  string `json:"tag"`
	Text *Text  `json:"text,omitempty"`
}

type Header struct {
	Title Text `json:"title"`
}

type WebhookMessageCard struct {
	Elements []*Element `json:"elements"`
	Header   *Header    `json:"header"`
}

type WebHookMessage struct {
	Timestamp string              `json:"timestamp"`
	Sign      string              `json:"sign"`
	MsgType   string              `json:"msg_type"`
	Card      *WebhookMessageCard `json:"card"`
}
