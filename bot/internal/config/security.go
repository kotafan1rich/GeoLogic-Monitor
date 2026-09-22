package config

type Security struct {
	MaxBotToken     string `env:"MAX_BOT_TOKEN,required"`
	BotServiceToken string `env:"BOT_SERVICE_TOKEN,required"`
}
