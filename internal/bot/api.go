package bot

import (
	"context"
	"fmt"

	"github.com/Dsmit05/party-day-bot/internal/logger"
	"github.com/Dsmit05/party-day-bot/internal/models"
	"github.com/pkg/errors"
)

func (p *PartyDayBot) SendMsgAllUser(ctx context.Context, msg string) error {
	userIDuserName := p.cache.GetUsers(ctx)

	var newErr error

	for id := range userIDuserName {
		if err := p.bot.SendText(id, msg); err != nil {
			newErr = err
			logger.Error("p.bot.SendText()", err)
		}
	}

	return newErr
}

func (p *PartyDayBot) SendMsg(ctx context.Context, user models.User, msg string) error {
	adminIDs := p.cache.GetAdmins(ctx)

	var newErr error

	newMsg := fmt.Sprintf("FirstName: %v\nLastName: %v\nUserName: %v\n%v",
		user.FirstName, user.LastName, user.UserName, msg)

	for _, id := range adminIDs {
		if err := p.bot.SendText(id, newMsg); err != nil {
			newErr = errors.Wrap(newErr, err.Error())
			logger.Error("p.bot.SendText", err)
		}
	}

	return newErr
}

func (p *PartyDayBot) SendPhoto(ctx context.Context, user models.User, url string) error {
	userIDuserName := p.cache.GetAdmins(ctx)

	var newErr error

	msg := fmt.Sprintf("FirstName: %v\nLastName: %v\nUserName: %v\n%v",
		user.FirstName, user.LastName, user.UserName, url)

	for _, id := range userIDuserName {
		if err := p.bot.SendPhotoFromURL(id, url); err != nil {
			newErr = errors.Wrap(newErr, err.Error())
			logger.Error("p.bot.SendPhotoFromUrl()", err)
		}

		if err := p.bot.SendText(id, msg); err != nil {
			newErr = err
		}
	}

	return newErr
}
