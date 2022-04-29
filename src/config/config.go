package config

import (
	"context"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	botToken string
}

func (c *Config) Init(ctx context.Context) error {
	log.Println("config: init")

	_ = godotenv.Load()

	err := c.loadBotToken()
	if err != nil {
		return err
	}

	log.Printf("config.init: complete")

	return nil
}

func (c *Config) loadBotToken() error {
	c.botToken = os.Getenv("BOT_TOKEN")

	return nil
}

func (c *Config) GetBotToken() string {
	return c.botToken
}

func NewConfig() *Config {
	return &Config{}
}
