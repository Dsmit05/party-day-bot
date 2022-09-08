package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Dsmit05/party-day-bot/internal/api"
	"github.com/Dsmit05/party-day-bot/internal/bot"
	"github.com/Dsmit05/party-day-bot/internal/cache"
	"github.com/Dsmit05/party-day-bot/internal/config"
	"github.com/Dsmit05/party-day-bot/internal/logger"
	"github.com/Dsmit05/party-day-bot/internal/repositories"
)

func main() {
	// Init logger
	if err := logger.InitLogger(false, "logs.json"); err != nil {
		panic(err)
	}

	// Init config service
	cfg, err := config.NewConfig()
	if err != nil {
		logger.Fatal("config.NewConfig()", err)
	}

	logger.Info("config.NewConfig", cfg.String())

	var db *repositories.PostgresRepository
	if cfg.IsDBUse() {
		// connect DB
		db, err = repositories.NewPostgresRepository(cfg)
		if err != nil {
			logger.Fatal("repositories.NewPostgresRepository()", err)
		}
		defer db.Close()
	}

	// create cache
	userCache, err := cache.NewUserCache(db, cfg)
	if err != nil {
		logger.Fatal("cache.NewUserCache()", err)
	}

	// init Bot
	partyBot, err := bot.NewPartyDayBot(cfg, userCache, db)
	if err != nil {
		logger.Fatal("bot.NewPartyDayBot()", err)
	}

	go partyBot.Start()

	// init api for Bot
	server := api.NewServer(cfg, partyBot)
	go server.StartGRPC()
	go server.StartREST()

	// Graceful Stop
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	partyBot.Stop()
	server.Stop()
}
