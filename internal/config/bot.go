package config

type Bot struct {
	Token string
}

func botConfig() Bot {
	return Bot{
		Token: getEnv("BOT_TOKEN"),
	}
}
