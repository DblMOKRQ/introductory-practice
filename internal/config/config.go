package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Storage `yaml:"storage" env-required:"true"`
	Rest    `yaml:"rest" env-required:"true"`
}

type Storage struct {
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	DBName   string `yaml:"dbname" env-required:"true"`
	Sslmode  string `yaml:"sslmode" env-default:"false"`
}
type Rest struct {
	Host string `yaml:"host" env-required:"true"`
	Port string `yaml:"port" env-required:"true"`
}

func MustLoad() *Config {
	os.Setenv("CONFIG_PATH", "/op/introductory-practice/config/config.yaml")
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		panic("CONFIG_PATH is not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("Config file does not exist")
	}
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config")
	}
	return &cfg

}
