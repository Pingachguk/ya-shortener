package config

import (
	"github.com/pingachguk/ya-shortener/internal/storage"
)

type config struct {
	App             string `env:"APP" envDefault:"localhost:8080"`
	Base            string `env:"BASE" envDefault:"http://localhost:8000"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN" envDefault:"postgres://postgres:postgres@localhost:5432/shortener"`
}

var Config config

func InitConfig() {
	if Config == (config{}) {
		parseFlags(&Config)

		if (Config.FileStoragePath != "") && (storage.GetStorage() == nil) {
			storage.NewFileStorage(Config.FileStoragePath)
		}

		if (Config.DatabaseDSN != "") && (storage.GetDatabase() == nil) {
			storage.InitDatabase(Config.DatabaseDSN)
		}
	}
}
