package config

type Rating struct {
	RecalcCron string        `env:"RATING_RECALC_CRON" envDefault:"0 3 1 * *"`
	ConfigPath string        `env:"RATING_CONFIG_PATH" envDefault:"rating.yml"`
	Formula    RatingFormula `env:"-"`
}
