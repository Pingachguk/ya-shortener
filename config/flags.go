package config

import (
	"flag"
	"os"
)

func parseFlags(cfg *Config) {
	flag.StringVar(&cfg.Host, "h", "127.0.0.1", "Host for server")
	flag.StringVar(&cfg.Port, "p", "8080", "Port for server")
	flag.Parse()

	if host := os.Getenv("HOST"); host != "" {
		cfg.Host = host
	}
	if port := os.Getenv("PORT"); port != "" {
		cfg.Port = port
	}
}
