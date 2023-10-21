package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env           string `yaml:"env" env-default:"local"`
	HttpServer    `yaml:"http_server"`
	Storages      `yaml:"storages"`
	ActiveStorage string `yaml:"active_storage" env-default:"sqlite"`
	Cache         `yaml:"cache"`
}

type HttpServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"3s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	CtxTimeout  time.Duration `yaml:"ctx_timeout" env-default:"8s"`
}

type Storages struct {
	SQLite SQLiteConfig `yaml:"sqlite"`
	Mongo  MongoConfig  `yaml:"mongo"`
}

type Cache struct {
	URL string `yaml:"url" env-default:"redis://localhost:6379/0"`
}

type SQLiteConfig struct {
	StoragePath string `yaml:"storage_path" env-required:"true"`
}

type MongoConfig struct {
	URI        string `yaml:"uri" env-required:"true"`
	Database   string `yaml:"database" env-default:"url-db"`
	Collection string `yaml:"collection" env-default:"urls"`
}

func MustLoadConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
