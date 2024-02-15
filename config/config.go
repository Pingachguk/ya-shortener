package config

import (
	"context"
	"github.com/caarlos0/env"
	"github.com/pingachguk/ya-shortener/internal/storage"
	"github.com/rs/zerolog/log"
)

type config struct {
	App              string `env:"APP" envDefault:"localhost:8080"`
	Base             string `env:"BASE" envDefault:"http://localhost:8000"`
	FileStoragePath  string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN      string `env:"DATABASE_DSN"`
	JWTKey           string `env:"JWT_KEY" envDefault:"notSecure"`
	JWTExpireSeconds int    `env:"JWT_EXPIRE_SECONDS" envDefault:"600"`
}

var Config config

func InitConfig() {
	if Config == (config{}) {
		if err := env.Parse(&Config); err != nil {
			log.Panic().Err(err).Msgf("bad parse env to config")
		}
		parseFlags(&Config)

		if storage.GetStorage() == nil {
			if Config.DatabaseDSN != "" {
				storage.InitDatabase(context.Background(), Config.DatabaseDSN)
			} else if Config.FileStoragePath != "" {
				storage.InitFileStorage(context.Background(), Config.FileStoragePath)
			} else {
				storage.InitMemoryStorage()
			}
		}
	}
}
