package config

import (
	"flag"
	"os"
)

func parseFlags(cfg *Config) {
	flag.StringVar(&cfg.App, "a", "localhost:8080", "App address for server")
	flag.StringVar(&cfg.Base, "b", "localhost:8080", "Base address for short URL")
	flag.Parse()

	if app := os.Getenv("APP"); app != "" {
		cfg.App = app
	}
	if base := os.Getenv("BASE"); base != "" {
		cfg.Base = base
	}
}
