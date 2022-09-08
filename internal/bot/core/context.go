package core

import (
	"bytes"
	"context"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ContextBot struct {
	Req tgbotapi.Update
	bot *Bot
	Ctx context.Context
}

func (r *ContextBot) SendText(chatID int64, text string) error {
	return r.bot.SendText(chatID, text)
}

func (r *ContextBot) SendPhoto(chatID int64, fileID string) error {
	return r.bot.SendPhoto(chatID, fileID)
}

func (r *ContextBot) GetLink(fileID string) (string, error) {
	file, err := r.bot.api.GetFile(tgbotapi.FileConfig{FileID: fileID})
	if err != nil {
		return "", err
	}

	link := file.Link(r.bot.config.Token)

	return link, nil
}

func (r *ContextBot) GetCommandArg() []string {
	msg := r.Req.Message.Text
	return strings.Split(msg, " ")[1:]
}

func (r *ContextBot) GetUserAndChatID() (*tgbotapi.User, int64) {
	user := r.Req.Message.From
	chatID := r.Req.Message.Chat.ID

	return user, chatID
}

func (r *ContextBot) GetRoutersInfoNonePrivate() string {
	var buf bytes.Buffer

	for _, v := range r.bot.routes {
		if v.Event.Form == Secret {
			continue
		}

		if !v.Event.Private {
			buf.WriteString(string(v.Event.Form))
			buf.WriteString(" ")
			buf.WriteString(v.Event.Command)
			buf.WriteString(" - ")
			buf.WriteString(v.Description)
			buf.WriteString("\n")
		}

	}

	return buf.String()
}

func (r *ContextBot) GetRoutersInfoAll() string {
	var buf bytes.Buffer

	for _, v := range r.bot.routes {
		if v.Event.Form == Secret {
			continue
		}

		buf.WriteString(string(v.Event.Form))
		buf.WriteString(" ")
		buf.WriteString(v.Event.Command)
		buf.WriteString(" - ")
		buf.WriteString(v.Description)
		buf.WriteString("\n")
	}

	return buf.String()
}
