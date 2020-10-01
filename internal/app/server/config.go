package server

type Config struct {
	BindAddr  string `toml:"bind_addr"`
	AppSecret string `toml:"app_secret"`
	AppHost   string `toml:"app_host"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr:  ":8080",
		AppSecret: "",
		AppHost:   "localhost",
	}
}
