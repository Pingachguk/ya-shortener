package config

import (
	"flag"
	"os"
)

func parseFlags(cfg *config) {
	flag.StringVar(&cfg.App, "a", "localhost:8080", "App address for server")
	flag.StringVar(&cfg.Base, "b", "http://localhost:8080", "Base address for short URL")
	flag.Parse()

	if app := os.Getenv("APP"); app != "" {
		cfg.App = app
	}
	if base := os.Getenv("BASE"); base != "" {
		cfg.Base = base
	}
}
