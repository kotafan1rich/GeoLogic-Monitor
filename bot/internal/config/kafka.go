package config

import "strings"

type Kafka struct {
	BrokersRaw string `env:"KAFKA_BROKERS" envDefault:"kafka:19092"`
	Topic      string `env:"KAFKA_TOPIC" envDefault:"notifications"`
	GroupID    string `env:"KAFKA_GROUP_ID" envDefault:"bot-notifications"`
}

func (c Kafka) Brokers() []string {
	return strings.Split(c.BrokersRaw, ",")
}
