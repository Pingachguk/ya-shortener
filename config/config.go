package config

import (
	"context"
	"github.com/pingachguk/ya-shortener/internal/storage"
)

type config struct {
	App             string `env:"APP" envDefault:"localhost:8080"`
	Base            string `env:"BASE" envDefault:"http://localhost:8000"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

var Config config

func InitConfig() {
	if Config == (config{}) {
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
