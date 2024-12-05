package service

import (
	"fmt"

	"github.com/HardDie/tg_bot_actions/internal/utils"
)

const (
	penisEmoji = "📏"
)

type PenisService struct {
	nouns []string
}

func NewPenisService() *PenisService {
	s := PenisService{}
	s.initNouns()
	return &s
}

func (s PenisService) GenerateSize() int {
	return utils.Random(40) + 1
}

func (s PenisService) GenerateDescription(size int) string {
	return fmt.Sprintf("%s %s у меня %dсм %s", penisEmoji, s.getNoun(), size, s.getSmile(size))
}

func (s PenisService) getSmile(size int) string {
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
		return sadSmiles[utils.Random(len(sadSmiles))]
	}
	return happySmiles[utils.Random(len(happySmiles))]
}
func (s PenisService) getNoun() string {
	return s.nouns[utils.Random(len(s.nouns))]
}

func (s *PenisService) initNouns() {
	nouns := map[string]struct{}{
		"Авторитет":         {},
		"Волшебная палочка": {},
		"Елда":              {},
		"Козырь в рукаве":   {},
		"Лоллипап":          {},
		"Лысый Джонни Синс": {},
		"Младший":           {},
		"Морковка":          {},
		"Мясная сигара":     {},
		"Нагибатель":        {},
		"Пенис":             {},
		"Песис":             {},
		"Писюлька":          {},
		"Питон в кустах":    {},
		"Стручок":           {},
		"Талант":            {},
		"Третья нога":       {},
		"Удав":              {},
		"Хоботок":           {},
		"Членохер":          {},
		"Чупачупс":          {},
		"Шишка":             {},
		"Шланг":             {},
	}
	s.nouns = make([]string, 0, len(nouns))
	for noun := range nouns {
		s.nouns = append(s.nouns, noun)
	}
}
