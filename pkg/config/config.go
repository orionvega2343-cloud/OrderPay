package config

import (
	"log"
	"log/slog"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

//Реализуем чтения .yml конфига из корня,
//в go syntax

type DB struct {
	Host     string `yaml:"host" required:"true"`
	Name     string `yaml:"name" required:"true"`
	Port     int    `yaml:"port" required:"true"`
	Username string `yaml:"user" required:"true"`
	Password string `env:"DB_PASS" required:"true"`
	SslMode  string `yaml:"ssl_mode" required:"true"`
}

type Server struct {
	Port string `yaml:"port" env-default:"8080"`
}

type Config struct {
	DB     DB     `yaml:"db"`
	Server Server `yaml:"server"`
}

//Функция подключения конфига

func MustLoad() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		slog.Info("Error loading .env file", "error", err)
	}

	cfg := Config{}

	err = cleanenv.ReadConfig("configs/config_dev.yml", &cfg)
	if err != nil {
		log.Fatal("Error loading config")
	}
	return &cfg
}
