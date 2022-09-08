package core

import (
	"context"
	"strings"
	"sync"

	"github.com/Dsmit05/party-day-bot/internal/consts"
	"github.com/Dsmit05/party-day-bot/internal/logger"
	"github.com/Dsmit05/party-day-bot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
)

type cacheI interface {
	CheckAccess(ctx context.Context, userID int64) bool
	AddUser(ctx context.Context, user models.User) error
	CheckUser(ctx context.Context, userID int64) bool
}

type Bot struct {
	api     *tgbotapi.BotAPI
	updates tgbotapi.UpdatesChannel
	cache   cacheI
	routes  Routes
	config  Config
}

func NewBot(config Config, cache cacheI) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return nil, errors.Wrap(err, "tgbotapi.NewBotAPI() error:")
	}

	api.Debug = config.Debug
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := api.GetUpdatesChan(u)

	return &Bot{
		api:     api,
		updates: updates,
		config:  config,
		cache:   cache,
	}, nil
}

func (b *Bot) Start() {
	b.routes.Log()

	wg := &sync.WaitGroup{}

	for i := 0; i < b.config.Worker; i++ {
		wg.Add(1)
		b.process(context.Background(), wg)
	}

	wg.Wait()
}

func (b *Bot) Stop() {
	b.api.StopReceivingUpdates()
}

func (b *Bot) AddRoute(event Event, handler HandlerFunc, description string) {
	b.routes = append(b.routes, Route{
		Event:       event,
		HandlerName: nameOfFunction(handler),
		Description: description,
		HandlerFunc: handler,
	})

}

func (b *Bot) SendText(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := b.api.Send(msg)

	return err
}

func (b *Bot) SendPhoto(chatID int64, fileID string) error {
	photo := tgbotapi.NewInputMediaPhoto(tgbotapi.FileID(fileID))
	mgc := tgbotapi.NewMediaGroup(chatID, []interface{}{photo})

	// Todo: telegram Bad Request error: type of file mismatch
	// When sending a document without compression from a PC, or the file itself, there is no forwarding to the Admin
	_, err := b.api.SendMediaGroup(mgc)

	return err
}

func (b *Bot) SendPhotoFromURL(chatID int64, url string) error {
	file := tgbotapi.FileURL(url)
	photo := tgbotapi.NewInputMediaPhoto(file)
	mgc := tgbotapi.NewMediaGroup(chatID, []interface{}{photo})

	_, err := b.api.SendMediaGroup(mgc)

	return err
}

// process perform the main work of message processing and sends the result to the user.
func (b *Bot) process(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case update, ok := <-b.updates:
			if !ok {
				logger.Info("b.updates", "channel updates close, stop process")
				return
			}

			form, err := b.checkEvent(update)
			if err != nil {
				logger.Error("checkEvent()", err)
				continue
			}

			if err := b.addUserInfo(ctx, update); err != nil {
				// тут может быть ошибка только с бд
				logger.Error("addUserInfo()", err)
			}

			if err := b.routing(ctx, form, update); err != nil {
				logger.Error("routing() error", err)
			}

		case <-ctx.Done():
			logger.Info("process ctx.Done()", ctx.Err().Error())
			return
		}
	}

}

func (b *Bot) routing(ctx context.Context, form *Event, update tgbotapi.Update) error {
	for i := 0; i < len(b.routes); i++ {
		route := b.routes[i]
		if route.Event.IsEqualRoute(form) {
			if !route.Event.Private {
				return route.HandlerFunc(ContextBot{
					Req: update,
					bot: b,
					Ctx: ctx,
				})
			}

			if b.cache.CheckAccess(ctx, update.Message.From.ID) {
				return route.HandlerFunc(ContextBot{
					Req: update,
					bot: b,
					Ctx: ctx,
				})
			}

			return errors.Errorf("user %v without access, try use handler: %v ",
				update.Message.From.ID, route.HandlerName)
		}
	}

	return errors.New("not equal route")
}

func (b *Bot) addUserInfo(ctx context.Context, update tgbotapi.Update) error {
	if b.cache.CheckUser(ctx, update.Message.From.ID) {
		return nil
	}

	user := models.User{
		TgID:      update.Message.From.ID,
		ChatID:    update.Message.Chat.ID,
		Role:      consts.RoleOrdinary,
		FirstName: update.Message.From.FirstName,
		LastName:  update.Message.From.LastName,
		UserName:  update.Message.From.UserName,
	}

	return b.cache.AddUser(ctx, user)
}

func (b *Bot) checkEvent(update tgbotapi.Update) (*Event, error) {
	if update.Message == nil {
		return nil, errors.New("unknown event")
	}

	photo := update.Message.Photo
	if photo != nil {
		return &Event{Form: Photo}, nil
	}

	document := update.Message.Document
	if document != nil {
		return &Event{Form: Document}, nil
	}

	text := update.Message.Text
	if text != "" {
		if text[0] == '/' {
			cmd := strings.Split(text, " ")[0]
			if cmd == b.config.SecretCMD {
				return &Event{
					Form: Secret,
				}, nil
			}

			return &Event{
				Form:    Command,
				Command: cmd,
			}, nil
		}

		if len(text) > 10 {
			return &Event{
				Form: Text,
			}, nil
		}

		return &Event{Form: Unknown}, nil
	}

	return nil, errors.New("unknown event")
}
