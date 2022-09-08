package media

import (
	"context"
	"fmt"

	"github.com/Dsmit05/party-day-bot/internal/bot/core"
	"github.com/Dsmit05/party-day-bot/internal/logger"
	"github.com/Dsmit05/party-day-bot/internal/models"
	"github.com/pkg/errors"
)

type cacheI interface {
	GetAdmins(ctx context.Context) []int64
}

type dbI interface {
	CreateFile(ctx context.Context, file models.File) error
}

type cfgI interface {
	IsDBUse() bool
}

func New(cache cacheI, db dbI, cfg cfgI) *media {
	return &media{cache, db, cfg}
}

type media struct {
	cacheI
	dbI
	cfgI
}

func (m *media) PhotoSendAdmin(c core.ContextBot) error {
	user, _ := c.GetUserAndChatID()

	foto := c.Req.Message.Photo
	maxFoto := foto[len(foto)-1].FileID

	link, err := c.GetLink(maxFoto)
	if err != nil {
		return err
	}

	logger.Info("New link", link)

	if m.IsDBUse() {
		if err := m.CreateFile(c.Ctx, models.File{
			TgID:     maxFoto,
			URL:      link,
			UserTgID: user.ID,
		}); err != nil {
			logger.Error("media.CreateFile()", err)
		}
	}

	admins := m.GetAdmins(c.Ctx)

	for _, v := range admins {
		if err := c.SendPhoto(v, maxFoto); err != nil {
			return err
		}

		msg := fmt.Sprintf("FirstName: %v\nLastName: %v\nUserName: %v\n%v",
			user.FirstName, user.LastName, user.UserName, link)

		if err := c.SendText(v, msg); err != nil {
			return err
		}
	}

	return nil
}

func (m *media) DocumentSendAdmin(c core.ContextBot) error {
	user, _ := c.GetUserAndChatID()
	file := c.Req.Message.Document.FileID

	link, err := c.GetLink(file)
	if err != nil {
		return err
	}

	logger.Info("New link", link)

	if m.IsDBUse() {
		fileInfo := models.File{
			TgID:     file,
			URL:      link,
			UserTgID: user.ID,
		}

		if err := m.CreateFile(c.Ctx, fileInfo); err != nil {
			logger.Error("media.CreateFile()", err)
		}
	}

	admins := m.GetAdmins(c.Ctx)

	for _, v := range admins {
		if err := c.SendPhoto(v, file); err != nil {
			return err
		}

		msg := fmt.Sprintf("FirstName: %v\nLastName: %v\nUserName: %v\n%v",
			user.FirstName, user.LastName, user.UserName, link)

		if err := c.SendText(v, msg); err != nil {
			return errors.Wrap(err, "с.SendText() err")
		}
	}

	return nil
}

func (m *media) PhotoInfo(c core.ContextBot) error {
	user, chatID := c.GetUserAndChatID()

	foto := c.Req.Message.Photo
	maxFoto := foto[len(foto)-1].FileID

	if err := c.SendPhoto(chatID, maxFoto); err != nil {
		return err
	}

	link, err := c.GetLink(maxFoto)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("FirstName: %v\nLastName: %v\nUserName: %v\n%v",
		user.FirstName, user.LastName, user.UserName, link)

	if err := c.SendText(chatID, msg); err != nil {
		return err
	}

	return nil
}

func (m *media) DocumentInfo(c core.ContextBot) error {
	user, chatID := c.GetUserAndChatID()

	file := c.Req.Message.Document.FileID

	if err := c.SendPhoto(chatID, file); err != nil {
		return err
	}

	link, err := c.GetLink(file)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("FirstName: %v\nLastName: %v\nUserName: %v\n%v",
		user.FirstName, user.LastName, user.UserName, link)

	if err := c.SendText(chatID, msg); err != nil {
		return err
	}

	return nil
}
