package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/kotstok/goblin/internal/app/server"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/chat.toml", "path to config file")
}

func main() {
	flag.Parse()

	config := server.NewConfig()
	_, err := toml.DecodeFile(configPath, config)

	if err != nil {
		log.Fatal(err)
	}

	s := server.New(config)

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
