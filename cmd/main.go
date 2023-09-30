package main

import (
	"golang.org/x/net/context"

	"github.com/HardDie/tg_bot_actions/internal/config"
	"github.com/HardDie/tg_bot_actions/internal/logger"
	"github.com/HardDie/tg_bot_actions/internal/repository"
	"github.com/HardDie/tg_bot_actions/internal/server"
	"github.com/HardDie/tg_bot_actions/internal/service"
)

func main() {
	cfg := config.Get()

	cache := repository.NewCacheRepository(cfg.Cache)
	penis := service.NewPenisService()
	gayMeter := service.NewGayMeterService()
	pressure := service.NewPressureService()
	criminal, err := service.NewCriminalService()
	if err != nil {
		logger.Error.Fatal(err)
	}
	pokemon, err := service.NewPokemonService()
	if err != nil {
		logger.Error.Fatal(err)
	}

	bot, err := server.NewTelegramServer(cfg.Bot, cache, penis, gayMeter, criminal, pokemon, pressure)
	if err != nil {
		logger.Error.Fatal(err)
	}

	// TODO: Add graceful shutdown
	logger.Info.Println("Waiting for messages...")
	err = bot.Run(context.Background())
	if err != nil {
		logger.Error.Fatal(err)
	}
}
