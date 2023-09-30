package server

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"golang.org/x/net/context"

	"github.com/HardDie/tg_bot_actions/internal/config"
	"github.com/HardDie/tg_bot_actions/internal/logger"
	"github.com/HardDie/tg_bot_actions/internal/models"
	"github.com/HardDie/tg_bot_actions/internal/repository"
	"github.com/HardDie/tg_bot_actions/internal/service"
)

type TelegramServer struct {
	bot *tgbotapi.BotAPI

	cache    *repository.CacheRepository
	penis    *service.PenisService
	gayMeter *service.GayMeterService
	criminal *service.CriminalService
	pokemon  *service.PokemonService
	pressure *service.PressureService
}

func NewTelegramServer(
	cfg config.Bot,
	cache *repository.CacheRepository,
	penis *service.PenisService,
	gayMeter *service.GayMeterService,
	criminal *service.CriminalService,
	pokemon *service.PokemonService,
	pressure *service.PressureService,
) (*TelegramServer, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		return nil, fmt.Errorf("error init telegram bot: %w", err)
	}

	return &TelegramServer{
		bot:      bot,
		cache:    cache,
		penis:    penis,
		gayMeter: gayMeter,
		criminal: criminal,
		pokemon:  pokemon,
		pressure: pressure,
	}, nil
}

func (s TelegramServer) Run(_ context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	for update := range s.bot.GetUpdatesChan(u) {
		if update.InlineQuery == nil {
			logger.Debug.Println("Skip non inline request")
			continue
		}

		var data *models.Cache

		// Check in cache and extract
		username := s.extractUsername(update.InlineQuery)
		if username != "" {
			data = s.cache.Get(username)
		} else {
			logger.Warn.Println("Username not found in inline query")
		}

		// Generate new
		if data == nil {
			data = &models.Cache{
				CockSize: s.penis.GenerateSize(),
				GayMeter: s.gayMeter.GenerateValue(),
				Criminal: s.criminal.GenerateCriminalIndex(),
				Pokemon:  s.pokemon.GeneratePokemonIndex(),
				Pressure: s.pressure.GeneratePressure(),
			}
		}
		// Update cache
		if username != "" {
			s.cache.Set(username, *data)
		}

		articleCock := s.newArticle("Петух размер", s.penis.GenerateDescription(data.CockSize), true)
		articleGay := s.newArticle("На сколько я гей?", s.gayMeter.GenerateDescription(data.GayMeter), true)
		articleCriminal := s.newArticle("Твоя статья УК РФ", s.criminal.GenerateDescription(data.Criminal), true)
		articlePressure := s.newArticle("Давление у меня?", s.pressure.GenerateDescription(data.Pressure), true)
		articlePokemon := s.newArticle("Это что за покемон?", s.pokemon.GenerateDescription(data.Pokemon), false)

		articleAllIn := s.newArticle("Все и сразу!",
			s.penis.GenerateDescription(data.CockSize)+"\n\n"+
				s.gayMeter.GenerateDescription(data.GayMeter)+"\n\n"+
				s.pressure.GenerateDescription(data.Pressure)+"\n\n"+
				s.criminal.GenerateDescription(data.Criminal),
			true,
		)

		s.send(update.InlineQuery.ID, articleCock, articleGay, articleCriminal, articlePressure, articlePokemon, articleAllIn)
	}
	return nil
}

func (s TelegramServer) send(id string, articles ...any) {
	inlineConf := tgbotapi.InlineConfig{
		InlineQueryID: id,
		IsPersonal:    true,
		CacheTime:     1,
		Results:       articles,
	}
	if _, err := s.bot.Send(inlineConf); err != nil {
		logger.Error.Println("send error:", err)
	}
}
func (s TelegramServer) extractUsername(m *tgbotapi.InlineQuery) string {
	if m.From == nil {
		return ""
	}
	if m.From.UserName != "" {
		return m.From.UserName
	}
	if m.From.ID == 0 {
		return ""
	}
	return fmt.Sprintf("%d", m.From.ID)
}
func (s TelegramServer) newArticle(title, text string, disableWebPreview bool) tgbotapi.InlineQueryResultArticle {
	article := tgbotapi.NewInlineQueryResultArticleHTML(
		uuid.New().String(),
		title,
		text,
	)
	val, ok := article.InputMessageContent.(tgbotapi.InputTextMessageContent)
	if !ok {
		logger.Error.Println("error cast InputMessageContent to type InputTextMessageContent")
		return article
	}

	val.DisableWebPagePreview = disableWebPreview
	article.InputMessageContent = val
	return article
}
