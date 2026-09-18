package config

import (
	"strconv"
	"time"
)

type Database struct {
	Host            string        `env:"POSTGRES_HOST" envDefault:"localhost"`
	Port            int           `env:"POSTGRES_PORT" envDefault:"5432"`
	User            string        `env:"POSTGRES_USER" envDefault:"postgres"`
	Password        string        `env:"POSTGRES_PASSWORD" envDefault:"postgres"`
	Name            string        `env:"POSTGRES_NAME" envDefault:"postgres"`
	SSLMode         string        `env:"POSTGRES_SSL_MODE" envDefault:"disable"`
	MinIdleConns    int           `env:"MIN_IDLE_CONNS" envDefault:"0"`
	MaxOpenConns    int           `env:"MAX_OPEN_CONNS" envDefault:"100"`
	MaxConnLifetime time.Duration `env:"MAX_CONN_LIFETIME" envDefault:"5m"`
}

func (d *Database) DSN() string {
	return "host=" + d.Host +
		" port=" + strconv.Itoa(d.Port) +
		" user=" + d.User +
		" password=" + d.Password +
		" dbname=" + d.Name +
		" sslmode=" + d.SSLMode
}
