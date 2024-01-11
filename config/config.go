package config

type config struct {
	App  string `env:"APP" envDefault:"localhost:8080"`
	Base string `env:"BASE" envDefault:"http://localhost:8000"`
}

var Config config

func InitConfig() {
	if Config == (config{}) {
		parseFlags(&Config)
	}
}
