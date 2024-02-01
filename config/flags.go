package config

import (
	"flag"
	"os"
)

func parseFlags(cfg *config) {
	flag.StringVar(&cfg.App, "a", "localhost:8080", "App address for server")
	flag.StringVar(&cfg.Base, "b", "http://localhost:8080", "Base address for short URL")
	flag.StringVar(&cfg.FileStoragePath, "f", "/tmp/short-url-db.json", "File storage for URLs. Example: /tmp/short-url-db.json")
	flag.StringVar(
		&cfg.DatabaseDSN,
		"d",
		"",
		"DSN for connect to database. Example: postgres://postgres:postgres@localhost:5432/shortener",
	)
	flag.Parse()

	if app := os.Getenv("APP"); app != "" {
		cfg.App = app
	}
	if base := os.Getenv("BASE"); base != "" {
		cfg.Base = base
	}
	if fileStoragePath := os.Getenv("FILE_STORAGE_PATH"); fileStoragePath != "" {
		cfg.FileStoragePath = fileStoragePath
	}
	if databaseDSN := os.Getenv("DATABASE_DSN"); databaseDSN != "" {
		cfg.DatabaseDSN = databaseDSN
	}
}
