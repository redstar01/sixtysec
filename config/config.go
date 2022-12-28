package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config - the main app config.
	Config struct {
		AppName                string `env-required:"true" env:"APP_NAME"`
		LogLevel               string `env-required:"true" env:"LOG_LEVEL"`
		TelegramToken          string `env-required:"true" env:"TELEGRAM_TOKEN"`
		GameSpeed              int    `env-required:"true" env:"GAME_SPEED"`
		GameQuestionsCount     int    `env-required:"true" env:"GAME_QUESTIONS_COUNT"`
		CacheDefaultExpiration int    `env-required:"true" env:"CACHE_DEFAULT_EXPIRATION"`
		CacheCleanupInterval   int    `env-required:"true" env:"CACHE_CLEANUP_INTERVAL"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
