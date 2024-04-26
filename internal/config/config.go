package config

type Config struct {
	SlackAPIToken string
	SlackChannel  string
}

func New() *Config {
	return &Config{}
}
