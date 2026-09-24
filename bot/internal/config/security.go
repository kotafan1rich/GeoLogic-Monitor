package config

type Security struct {
	WebhookUrl      string `env:"WEBHOOK_URL"`
	WebhookSecret   string `env:"WEBHOOK_SECRET"`
	MaxBotToken     string `env:"MAX_BOT_TOKEN,required"`
	BotServiceToken string `env:"BOT_SERVICE_TOKEN,required"`
}
