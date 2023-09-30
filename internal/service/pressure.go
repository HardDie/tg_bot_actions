package service

import (
	"fmt"
	"math"

	"github.com/HardDie/tg_bot_actions/internal/utils"
)

const (
	LowestSys  = 80
	HighestSys = 190

	LowestDiaPerc  = 60
	HighestDiaPerc = 85
)

type PressureService struct {
	prefixes []string
}

func NewPressureService() *PressureService {
	s := PressureService{}
	s.initPrefixes()
	return &s
}

func (s PressureService) GeneratePressure() string {
	sys, dia := s.generatePressureValues()
	return fmt.Sprintf("%d/%d", sys, dia)
}

func (s PressureService) GenerateDescription(pressure string) string {
	return fmt.Sprintf("%s давление у меня %s", s.getPrefix(), pressure)
}

func (s PressureService) getPrefix() string {
	return s.prefixes[utils.Random(len(s.prefixes))]
}

func (s PressureService) generatePressureValues() (int, int) {
	sysPressure := utils.Random(HighestSys-LowestSys+1) + LowestSys
	diaPressure := utils.Random(HighestDiaPerc-LowestDiaPerc+1) + LowestDiaPerc
	diaPressure = int(math.Floor(float64(sysPressure*diaPressure) / 100))
	return sysPressure, diaPressure
}

func (s *PressureService) initPrefixes() {
	prefixes := map[string]struct{}{
		"В любое время дня и ночи":                {},
		"Внезапно,":                               {},
		"Вы тут со своим этим, а":                 {},
		"Ёпсель-мопсель,":                         {},
		"Из-за погоды":                            {},
		"Можно конечно меряться сантиметрами, но": {},
		"Не знаю, как там у вас, а":               {},
		"Не знаю, нормально это или нет, но":      {},
		"Не то чтобы это было смертельно, но":     {},
		"Не хочу хвастаться, но":                  {},
		"Не хухры-мухры,":                         {},
		"Обычное дело,":                           {},
		"От стресса":                              {},
		"От таких новостей":                       {},
		"Офигеть,":                                {},
		"По жизни":                                {},
		"Пока не хлопну рюмочку,":                 {},
		"После восьми чашек кофе":                 {},
		"После трёх ступенек":                     {},
		"Последнее время":                         {},
		"Прочитал, что вы тут пишете, и теперь":   {},
		"С вашими игрушками":                      {},
		"Сколько уже не употребляю, а":            {},
		"Со вчерашнего дня":                       {},
		"Уже такой возраст, что":                  {},
		"Хоть стой, хоть падай, а":                {},
		"Целый день":                              {},
		"Чтоб вы знали,":                          {},
		"Чувствую себя отлично,":                  {},
		"Шутки шутками, а":                        {},
	}
	s.prefixes = make([]string, 0, len(prefixes))
	for prefixes := range prefixes {
		s.prefixes = append(s.prefixes, prefixes)
	}
}
