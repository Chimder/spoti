package config

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type EnvVars struct {
	Username       string
	Password       string
	SteamID        string
	SharedSecret   string
	IdentitySecret string
	DeviceID       string
	PostgresUrl    string
	Debug          bool
	Env            string
	LogLevel       string
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
			log.Warn().Msg(".env not found")
		}
		config = &EnvVars{
			Username:       setEnv("USERNAME", ""),
			Password:       setEnv("PASSWORD", ""),
			SteamID:        setEnv("STEAM_ID", ""),
			SharedSecret:   setEnv("SHARED_SECRET", ""),
			IdentitySecret: setEnv("IDENTITY_SECRET", ""),
			DeviceID:       setEnv("DEVICE_ID", ""),
			PostgresUrl:    setEnv("POSTGRES_URL", ""),
			Env:            setEnv("ENV", "dev"),
			LogLevel:       setEnv("LOG_LEVEL", "debug"),
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
