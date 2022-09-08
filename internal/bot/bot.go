package bot

import (
	"github.com/Dsmit05/party-day-bot/internal/bot/controllers/admin"
	"github.com/Dsmit05/party-day-bot/internal/bot/controllers/info"
	"github.com/Dsmit05/party-day-bot/internal/bot/controllers/media"
	"github.com/Dsmit05/party-day-bot/internal/bot/controllers/text"
	"github.com/Dsmit05/party-day-bot/internal/bot/core"
	"github.com/Dsmit05/party-day-bot/internal/config"
	"github.com/Dsmit05/party-day-bot/internal/logger"
)

type PartyDayBot struct {
	bot   *core.Bot
	cache cacheI
	db    dbI
	cfg   *config.Config
}

func NewPartyDayBot(cfg *config.Config, cache cacheI, db dbI) (*PartyDayBot, error) {
	coreCfg := core.Config{
		Token:     cfg.Bot.Key,
		Worker:    4,
		Timeout:   60,
		Debug:     cfg.Bot.Debug,
		SecretCMD: cfg.GetSecretCommand(),
	}

	bot, err := core.NewBot(coreCfg, cache)
	if err != nil {
		logger.Error("core.NewBot", err)
		return nil, err
	}

	routesMedia := media.New(cache, db, cfg)
	{
		bot.AddRoute(
			core.Event{Form: core.Photo},
			routesMedia.PhotoSendAdmin,
			"Send me a photo and I will save it")

		bot.AddRoute(
			core.Event{Form: core.Document},
			routesMedia.DocumentSendAdmin,
			"You can send an already saved photo or document")
	}

	routesInfo := info.New()
	{
		bot.AddRoute(core.Event{
			Form:    core.Command,
			Command: "/whoAmI",
			Private: false,
		}, routesInfo.WhoAmI, "Information about me")

		bot.AddRoute(core.Event{
			Form:    core.Command,
			Command: "/help",
			Private: false,
		}, routesInfo.HelpUser, "Information about the main teams")

		bot.AddRoute(core.Event{
			Form:    core.Command,
			Command: "/helpA",
			Private: true,
		}, routesInfo.HelpAdmin, "Information about the main Admin commands")
	}

	routesAdmin := admin.New(cache)
	{
		bot.AddRoute(core.Event{
			Form:    core.Command,
			Command: "/list",
			Private: true,
		}, routesAdmin.ListUsers, "Guest list")

		bot.AddRoute(core.Event{
			Form:    core.Command,
			Command: "/root",
			Private: true,
		}, routesAdmin.AddUserAccess, "Give the user rights, /root id")

		bot.AddRoute(core.Event{
			Form:    core.Command,
			Command: "/sendAll",
			Private: true,
		}, routesAdmin.SendMsgAllUsers, "send a message all guests, /sendAll your message")

		bot.AddRoute(core.Event{
			Form:    core.Secret,
			Private: false,
		}, routesAdmin.AddRoot, "Gives the user who knows the secret administrator rights")
	}

	routesText := text.New(cache, db, cfg)
	{
		bot.AddRoute(core.Event{
			Form: core.Text,
		}, routesText.TextSendAdmin, "You can send me your wishes or congratulations")
	}

	return &PartyDayBot{bot: bot, cache: cache, db: db, cfg: cfg}, nil
}

func (p *PartyDayBot) Start() {
	p.bot.Start()
}

func (p *PartyDayBot) Stop() {
	p.bot.Stop()
}
