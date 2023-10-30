package utils

import (
	crand "crypto/rand"
	"math/big"
	"math/rand"

	"github.com/HardDie/tg_bot_actions/internal/logger"
)

// Random generate random value from [0, value)
func Random(value int) int {
	bigValue, err := crand.Int(crand.Reader, big.NewInt(int64(value)))
	if err == nil {
		return int(bigValue.Int64())
	}
	logger.Error.Println("error generating crypto/rand value:", err.Error())
	return rand.Intn(value)
}
