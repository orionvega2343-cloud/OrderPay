package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

//Реализуем чтения .yml конфига из корня,
//в go syntax

type DB struct {
	Host     string `yaml:"host" required:"true"`
	Name     string `yaml:"name" required:"true"`
	Port     string `yaml:"port" required:"true"`
	Username string `yaml:"user" required:"true"`
	Password string `env:"DB_PASS" required:"true"`
	SslMode  string `yaml:"ssl_mode" required:"true"`
}

type Config struct {
	DB DB `yaml:"db"`
}

//Функция подключения конфига

func MustLoad() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := Config{}

	err = cleanenv.ReadConfig("config/config.yml", &cfg)
	if err != nil {
		log.Fatal("Error loading config")
	}
	return &cfg
}
