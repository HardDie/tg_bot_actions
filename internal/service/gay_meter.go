package service

import (
	"fmt"

	"github.com/HardDie/tg_bot_actions/internal/utils"
)

const (
	gayFlag = "🏳️‍🌈"
)

type GayMeterService struct {
}

func NewGayMeterService() *GayMeterService {
	return &GayMeterService{}
}

func (s GayMeterService) GenerateValue() int {
	return utils.Random(101)
}

func (s GayMeterService) GenerateDescription(value int) string {
	return fmt.Sprintf("%s ‍Я на %d%% гей!", gayFlag, value)
}
