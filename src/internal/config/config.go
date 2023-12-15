package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env			string		`yaml:"env" env-default:"local"`
	StoragePath	string		`yaml:"storage_path" env-required:"true"`
	HTTPServer  			`yaml:"http_server"`
}

type HTTPServer struct {
	Address			string			`yaml:"address" env-default:"localhost:8080"`
	Timeout			time.Duration	`yaml:"timeout" env-default:"4s"`
	IdleTimeout		time.Duration	`yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		log.Fatalf("config file is not exist: %s", cfgPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		log.Fatalf("error while reading config: %s", err)
	}

	return &cfg
}