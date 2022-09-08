package text

import (
	"context"
	"fmt"

	"github.com/Dsmit05/party-day-bot/internal/bot/core"
	"github.com/Dsmit05/party-day-bot/internal/logger"
)

type cacheI interface {
	GetAdmins(ctx context.Context) []int64
}

type dbI interface {
	CreateMessage(ctx context.Context, userTgID int64, text string) error
}

type cfgI interface {
	IsDBUse() bool
}

func New(cache cacheI, db dbI, cfg cfgI) *text {
	return &text{cache, db, cfg}
}

type text struct {
	cacheI
	dbI
	cfgI
}

func (t *text) TextSendAdmin(c core.ContextBot) error {
	user, _ := c.GetUserAndChatID()
	userText := c.Req.Message.Text
	admins := t.GetAdmins(c.Ctx)

	if t.IsDBUse() {
		if err := t.CreateMessage(c.Ctx, user.ID, userText); err != nil {
			logger.Error("text.CreateMessage()", err)
		}
	}

	for _, v := range admins {
		msg := fmt.Sprintf("FirstName: %v\nLastName: %v\nUserName: %v\n%v",
			user.FirstName, user.LastName, user.UserName, userText)

		if err := c.SendText(v, msg); err != nil {
			return err
		}
	}

	return nil
}
