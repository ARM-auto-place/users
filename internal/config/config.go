package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/rs/zerolog"
)

type PostgresConfig struct {
	Addr            string        `env:"PG_ADDR,required"`
	DB              string        `env:"PG_DB,required"`
	User            string        `env:"PG_USER,required"`
	Password        string        `env:"PG_PASSWORD"`
	MaxConn         int32         `env:"PG_MAX_CONN" envDefault:"10"`
	MaxIdleLifetime time.Duration `env:"PG_MAX_IDLE_LIFETIME" envDefault:"30s"`
	MaxLifetime     time.Duration `env:"PG_MAX_LIFETIME" envDefault:"3m"`
	ReplicaAddr     string        `env:"PG_REPLICA_ADDR"`
	UserRO          string        `env:"PG_REPLICA_USER"`
	PasswordRO      string        `env:"PG_REPLICA_PASSWORD"`
}

type Config struct {
	ServiceName    string   `env:"SERVICE_NAME" envDefault:"users"`
	LogLevel       string   `env:"LOGGER_LEVEL" envDefault:"info"`
	LogFormat      string   `env:"LOGGER_FORMAT"`
	ServiceBind    string   `env:"BIND_ADDR" envDefault:":8083"`
	HealthBind     string   `env:"BIND_HEALTH" envDefault:":9091"`
	AccessAddr     string   `env:"ACCESS_ADDRES,required"`
	AllowedOrigins []string `env:"ALLOWED_ORIGINS" envSeparator:","`

	MaxRequestBodySize   int           `env:"MAX_REQUEST_BODY_SIZE" envDefault:"16384"`
	MaxRequestHeaderSize int           `env:"MAX_REQUEST_HEADER_SIZE" envDefault:"16384"`
	ReadTimeout          time.Duration `env:"READ_TIMEOUT" envDefault:"15s"`
	WriteTimeout         time.Duration `env:"WRITE_TIMEOUT" envDefault:"15s"`
	ShutdownTimeout      time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"10s"`

	Users PostgresConfig
	Cars  PostgresConfig `envPrefix:"CARS_"`
}

func Values() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}
	accessURL, err := url.Parse(cfg.AccessAddr)
	if err != nil || accessURL.Hostname() == "" || (accessURL.Scheme != "http" && accessURL.Scheme != "https") {
		return Config{}, fmt.Errorf("ACCESS_ADDRES must be an absolute HTTP or HTTPS URL")
	}
	for i := range cfg.AllowedOrigins {
		cfg.AllowedOrigins[i] = strings.TrimSpace(cfg.AllowedOrigins[i])
	}
	return cfg, nil
}

func (cfg Config) Logger() zerolog.Logger {
	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	var logger zerolog.Logger
	if cfg.LogFormat == "console" {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	} else {
		logger = zerolog.New(os.Stdout)
	}
	return logger.Level(level).With().Timestamp().Str("service", cfg.ServiceName).Logger()
}
