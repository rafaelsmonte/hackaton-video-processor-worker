package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct{}

func NewConfig() Config {
	godotenv.Load()
	return Config{}
}
func (c *Config) GetConfig(name string) string {
	return os.Getenv(name)
}
