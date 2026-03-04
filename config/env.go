package config

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type EnvVars struct {
	PostgresUrl string
	Debug       bool
	Env         string
	LogLevel    string
}

var config *EnvVars
var once sync.Once

func GetEnv() *EnvVars {
	if config == nil {
		log.Fatal().Msg("LoadEnv didnt happen")
	}

	return config
}

func LoadEnv() *EnvVars {
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			log.Warn().Err(err).Msg(".env not found")
		} else {
			log.Info().Msg(".env loaded")
		}
		config = &EnvVars{
			PostgresUrl: setEnv("POSTGRES_URL", ""),
			Env:         setEnv("ENV", "dev"),
			LogLevel:    setEnv("LOG_LEVEL", "debug"),
		}
	})

	return config
}

func setEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	log.Info().Str("use default", defaultVal).Str("for key", key).Msg("ENV")
	return defaultVal
}
