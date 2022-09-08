package info

import (
	"bytes"
	"fmt"

	"github.com/Dsmit05/party-day-bot/internal/bot/core"
)

func New() *info {
	return &info{}
}

type info struct {
}

func (m *info) WhoAmI(c core.ContextBot) error {
	user, chatID := c.GetUserAndChatID()
	msg := fmt.Sprintf("UserName: %v\nChatID: %v\nUserID: %v\n",
		user.UserName, chatID, user.ID)

	return c.SendText(chatID, msg)
}

func (m *info) HelpUser(c core.ContextBot) error {
	_, chatID := c.GetUserAndChatID()
	cmd := c.GetRoutersInfoNonePrivate()

	var buf bytes.Buffer

	buf.WriteString("You can interact with me in the following ways:\n\n")
	buf.WriteString(cmd)

	return c.SendText(chatID, buf.String())
}

func (m *info) HelpAdmin(c core.ContextBot) error {
	_, chatID := c.GetUserAndChatID()
	cmd := c.GetRoutersInfoAll()

	var buf bytes.Buffer

	buf.WriteString("You can interact with me in the following ways:\n\n")
	buf.WriteString(cmd)

	return c.SendText(chatID, buf.String())
}
