package config

type Config struct {
	Host string `env:"HOST" envDefault:"127.0.0.1"`
	Port string `env:"PORT" envDefault:"8080"`
}

func New() Config {
	cfg := Config{}
	parseFlags(&cfg)

	return cfg
}
