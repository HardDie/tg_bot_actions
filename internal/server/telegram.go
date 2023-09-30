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
}

func NewTelegramServer(
	cfg config.Bot,
	cache *repository.CacheRepository,
	penis *service.PenisService,
	gayMeter *service.GayMeterService,
	criminal *service.CriminalService,
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
			}
		}
		// Update cache
		if username != "" {
			s.cache.Set(username, *data)
		}

		articleCock := s.newArticle("Петух размер", s.penis.GenerateDescription(data.CockSize))
		articleGay := s.newArticle("На сколько я гей?", s.gayMeter.GenerateDescription(data.GayMeter))
		articleCriminal := s.newArticle("Твоя статья УК РФ", s.criminal.GenerateDescription(data.Criminal))
		articleAllIn := s.newArticle("Все и сразу!",
			s.penis.GenerateDescription(data.CockSize)+"\n\n"+
				s.gayMeter.GenerateDescription(data.GayMeter)+"\n\n"+
				s.criminal.GenerateDescription(data.Criminal),
		)

		s.send(articleCock, articleGay, articleCriminal, articleAllIn)
	}
	return nil
}

func (s TelegramServer) send(articles ...any) {
	inlineConf := tgbotapi.InlineConfig{
		InlineQueryID: uuid.New().String(),
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
func (s TelegramServer) newArticle(title, text string) tgbotapi.InlineQueryResultArticle {
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

	val.DisableWebPagePreview = true
	article.InputMessageContent = val
	return article
}
