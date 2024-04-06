package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Port  string
	PgURL string
}

func NewConfig() *Config {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println(fmt.Errorf("error config file: %w", err))
		viper.SetDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432")
		viper.SetDefault("PORT", "8080")
	}

	return &Config{
		PgURL: viper.GetString("DATABASE_URL"),
		Port:  viper.GetString("PORT"),
	}
}
