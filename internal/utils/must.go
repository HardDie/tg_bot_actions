package utils

import "github.com/HardDie/tg_bot_actions/internal/logger"

func Must[T any](val T, err error) T {
	if err != nil {
		logger.Error.Println("function must be without error:", err.Error())
		panic("must be without error")
	}
	return val
}
