package config

type Config struct {
	App  string `env:"APP" envDefault:"localhost:8080"`
	Base string `env:"BASE" envDefault:"localhost:8000"`
}

func New() Config {
	cfg := Config{}
	parseFlags(&cfg)

	return cfg
}
