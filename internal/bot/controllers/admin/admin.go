package admin

import (
	"bytes"
	"context"
	"strconv"

	"github.com/Dsmit05/party-day-bot/internal/bot/core"
	"github.com/Dsmit05/party-day-bot/internal/consts"
)

type cacheI interface {
	UpdateRole(ctx context.Context, userID int64, role string) error
	GetUsers(ctx context.Context) map[int64]string
}

func New(cache cacheI) *admin {
	return &admin{cache: cache}
}

type admin struct {
	cache cacheI
}

func (a *admin) ListUsers(c core.ContextBot) error {
	_, chatID := c.GetUserAndChatID()

	var buf bytes.Buffer

	buf.WriteString("List of users: \n")
	buf.WriteString("ID - User Name \n")
	buf.WriteString("______________ \n")

	for i, v := range a.cache.GetUsers(c.Ctx) {
		buf.WriteString(strconv.Itoa(int(i)))
		buf.WriteString(" - ")
		buf.WriteString(v)
		buf.WriteString("\n")
	}

	return c.SendText(chatID, buf.String())
}

func (a *admin) AddUserAccess(c core.ContextBot) error {
	_, chatID := c.GetUserAndChatID()
	arg := c.GetCommandArg()

	if len(arg) == 0 {
		return c.SendText(chatID, "Add a user ID")
	}

	id, err := strconv.ParseInt(arg[0], 10, 0)
	if err != nil {
		return c.SendText(chatID, "Invalid user ID format")
	}

	if err := a.cache.UpdateRole(c.Ctx, id, consts.RoleAdmin); err != nil {
		return c.SendText(chatID, err.Error())
	}

	return c.SendText(chatID, "User got the rights")
}

func (a *admin) SendMsgAllUsers(c core.ContextBot) error {
	_, chatID := c.GetUserAndChatID()
	userIDuserName := a.cache.GetUsers(c.Ctx)
	arg := c.GetCommandArg()

	if len(arg) == 0 {
		return c.SendText(chatID, "Add the message body")
	}

	var buf bytes.Buffer
	for _, v := range arg {
		buf.WriteString(v)
		buf.WriteString(" ")
	}

	for id := range userIDuserName {
		if err := c.SendText(id, buf.String()); err != nil {
			return err
		}
	}

	return nil
}

func (a *admin) AddRoot(c core.ContextBot) error {
	user, chatID := c.GetUserAndChatID()

	if err := a.cache.UpdateRole(c.Ctx, user.ID, consts.RoleAdmin); err != nil {
		return c.SendText(chatID, err.Error())
	}

	return c.SendText(chatID, "You got the rights")
}
