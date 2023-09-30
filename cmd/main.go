package main

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"

	"github.com/HardDie/tg_bot_actions/internal/config"
	"github.com/HardDie/tg_bot_actions/internal/logger"
	"github.com/HardDie/tg_bot_actions/internal/models"
	"github.com/HardDie/tg_bot_actions/internal/repository"
	"github.com/HardDie/tg_bot_actions/internal/service"
)

func getUsername(m *tgbotapi.InlineQuery) *string {
	if m.From == nil {
		return nil
	}
	if m.From.UserName != "" {
		return &m.From.UserName
	}
	if m.From.ID == 0 {
		return nil
	}
	res := fmt.Sprintf("%d", m.From.ID)
	return &res
}

func main() {
	cfg := config.Get()

	// Создаем бота
	bot, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		logger.Error.Fatal(err)
	}

	// Подписываемся на обновления
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	cache := repository.NewCacheRepository(cfg.Cache)
	penis := service.NewPenisService()
	gayMeter := service.NewGayMeterService()
	criminalService, err := service.NewCriminalService()
	if err != nil {
		logger.Error.Fatal(err)
	}

	// В цикле обрабатываем события
	logger.Info.Println("Wait for messages...")
	for update := range bot.GetUpdatesChan(u) {
		if update.InlineQuery == nil {
			log.Println("Skip request", update.InlineQuery)
			continue
		}

		var data *models.Cache

		// Check in cache and extract
		username := getUsername(update.InlineQuery)
		if username != nil {
			data = cache.Get(*username)
		} else {
			logger.Warn.Println("Username not found in inline query")
		}

		// Generate new
		if data == nil {
			data = &models.Cache{
				CockSize: penis.GenerateSize(),
				GayMeter: gayMeter.GenerateValue(),
				Criminal: criminalService.GenerateCriminalIndex(),
			}
		}
		// Update cache
		if username != nil {
			cache.Set(*username, *data)
		}

		articleCock := tgbotapi.NewInlineQueryResultArticleHTML(
			uuid.New().String(),
			"Петух размер",
			penis.GenerateDescription(data.CockSize),
		)

		articleGay := tgbotapi.NewInlineQueryResultArticleHTML(
			uuid.New().String(),
			"На сколько я гей?",
			gayMeter.GenerateDescription(data.GayMeter),
		)

		articleCriminal := tgbotapi.NewInlineQueryResultArticleHTML(
			uuid.New().String(),
			"Твоя статья УК РФ",
			criminalService.GenerateDescription(data.Criminal),
		)
		val, ok := articleCriminal.InputMessageContent.(tgbotapi.InputTextMessageContent)
		if ok {
			val.DisableWebPagePreview = true
			articleCriminal.InputMessageContent = val
		}

		articleAllIn := tgbotapi.NewInlineQueryResultArticleHTML(
			uuid.New().String(),
			"Все и сразу!",
			penis.GenerateDescription(data.CockSize)+"\n\n"+
				gayMeter.GenerateDescription(data.GayMeter)+"\n\n"+
				criminalService.GenerateDescription(data.Criminal),
		)
		val, ok = articleAllIn.InputMessageContent.(tgbotapi.InputTextMessageContent)
		if ok {
			val.DisableWebPagePreview = true
			articleAllIn.InputMessageContent = val
		}

		inlineConf := tgbotapi.InlineConfig{
			InlineQueryID: update.InlineQuery.ID,
			IsPersonal:    true,
			CacheTime:     1,
			Results:       []interface{}{articleCock, articleGay, articleCriminal, articleAllIn},
		}

		if _, err := bot.Send(inlineConf); err != nil {
			logger.Error.Println("send error:", err)
		}
	}
}
