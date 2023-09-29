package main

import (
	crand "crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/jellydator/ttlcache/v3"

	"github.com/HardDie/tg_bot_actions/internal/config"
	"github.com/HardDie/tg_bot_actions/internal/logger"
	"github.com/HardDie/tg_bot_actions/internal/models"
)

const (
	gayFlag      = "🏳️‍🌈"
	criminalBook = "📕"
)

type Cache struct {
	CockSize int `json:"cockSize"`
	GayMeter int `json:"gayMeter"`
	Criminal int `json:"criminal"`
}

func random(value int) int {
	bigValue, err := crand.Int(crand.Reader, big.NewInt(int64(value)))
	if err == nil {
		return int(bigValue.Int64())
	}
	logger.Error.Println("error generating crypto/rand value:", err.Error())
	return rand.Intn(value)
}
func getSmile(size int) string {
	sadSmiles := []string{
		"😞",
		"😨",
		"😒",
		"😣",
	}

	happySmiles := []string{
		"😏",
		"😁",
		"😱",
		"😯",
	}
	if size < 17 {
		return sadSmiles[random(3)]
	}
	return happySmiles[random(3)]
}
func getNoun() string {
	nouns := map[string]struct{}{
		"Талант":            {},
		"Шишка":             {},
		"Козырь в рукаве":   {},
		"Шланг":             {},
		"Чупачупс":          {},
		"Лысый Джонни Синс": {},
		"Волшебная палочка": {},
		"Лоллипап":          {},
		"Пенис":             {},
		"Нагибатель":        {},
		"Младший":           {},
		"Песис":             {},
		"Елда":              {},
		"Третья нога":       {},
		"Писюлька":          {},
		"Хоботок":           {},
		"Стручок":           {},
		"Авторитет":         {},
		"Удав":              {},
		"Морковка":          {},
		"Мясная сигара":     {},
		"Членохер":          {},
		"Питон в кустах":    {},
	}
	nounsList := make([]string, 0, len(nouns))
	for noun := range nouns {
		nounsList = append(nounsList, noun)
	}
	return nounsList[random(len(nouns))]
}

func newCriminals() []models.Criminal {
	file, err := os.Open("criminals.json")
	if err != nil {
		logger.Error.Fatal("error open criminals.json file", err.Error())
	}
	defer file.Close()

	var criminals []models.Criminal
	err = json.NewDecoder(file).Decode(&criminals)
	if err != nil {
		logger.Error.Fatal("error parse criminals.json file", err.Error())
	}

	return criminals
}

func getUsername(m *tgbotapi.InlineQuery) *string {
	if m.From == nil {
		return nil
	}
	if m.From.UserName == "" {
		return nil
	}
	return &m.From.UserName
}

func main() {
	cfg := config.Get()
	cache := ttlcache.New[string, Cache](
		ttlcache.WithTTL[string, Cache](cfg.Cache.Period),
	)
	criminals := newCriminals()

	// Создаем бота
	bot, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		logger.Error.Fatal(err)
	}

	// Подписываемся на обновления
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// В цикле обрабатываем события
	logger.Info.Println("Wait for messages...")
	for update := range bot.GetUpdatesChan(u) {
		if update.InlineQuery == nil {
			log.Println("Skip request", update.InlineQuery)
			continue
		}

		var data *Cache

		// Check in cache and extract
		username := getUsername(update.InlineQuery)
		if username == nil {
			logger.Warn.Println("Username not found in inline query")
		}
		if username != nil {
			item := cache.Get(*username)
			if item != nil {
				val := item.Value()
				data = &val
			}
		}

		// Generate new
		if data == nil {
			data = &Cache{
				CockSize: random(39) + 1,
				GayMeter: random(101),
				Criminal: random(len(criminals)),
			}
		}
		// Update cache
		if username != nil {
			cache.Set(*username, *data, ttlcache.DefaultTTL)
		}

		cockSmile := getSmile(data.CockSize)
		cockNoun := getNoun()
		articleCock := tgbotapi.NewInlineQueryResultArticleHTML(
			uuid.New().String(),
			"Петух размер",
			fmt.Sprintf("%s у меня %dсм %s", cockNoun, data.CockSize, cockSmile),
		)
		//articleCock.Description = "Петух размер"

		articleGay := tgbotapi.NewInlineQueryResultArticleHTML(
			uuid.New().String(),
			"На сколько я гей?",
			fmt.Sprintf("%s ‍Я на %d%% гей!", gayFlag, data.GayMeter),
		)
		//articleGay.Description = "На сколько я гей"

		articleCriminal := tgbotapi.NewInlineQueryResultArticleHTML(
			uuid.New().String(),
			"Твоя статья УК РФ",
			fmt.Sprintf(`%s <u>Твоя статья УК РФ</u>:
<a href="%s"><b>%s</b></a> - %s`, criminalBook, criminals[data.Criminal].Link, criminals[data.Criminal].Number, criminals[data.Criminal].Description),
		)
		val, ok := articleCriminal.InputMessageContent.(tgbotapi.InputTextMessageContent)
		if ok {
			val.DisableWebPagePreview = true
			articleCriminal.InputMessageContent = val
		}
		//articleCriminal.Description = "Твоя статья УК РФ"

		articleAllIn := tgbotapi.NewInlineQueryResultArticleHTML(
			uuid.New().String(),
			"Все и сразу!",
			fmt.Sprintf("%s у меня %dсм %s", cockNoun, data.CockSize, cockSmile)+
				"\n\n"+
				fmt.Sprintf("%s ‍Я на %d%% гей!", gayFlag, data.GayMeter)+
				"\n\n"+
				fmt.Sprintf(`%s <u>Твоя статья УК РФ</u>:
<a href="%s"><b>%s</b></a> - %s`, criminalBook, criminals[data.Criminal].Link, criminals[data.Criminal].Number, criminals[data.Criminal].Description),
		)
		val, ok = articleAllIn.InputMessageContent.(tgbotapi.InputTextMessageContent)
		if ok {
			val.DisableWebPagePreview = true
			articleAllIn.InputMessageContent = val
		}
		//articleAllIn.Description = "Все и сразу!"

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
