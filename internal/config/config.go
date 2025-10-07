package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DB DBConfig
}

type DBConfig struct {
	Host     string
	Password string
	User     string
	DBName   string
	Port     string
	SSLMode  string
}

var AppConfig *Config

func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("⚠️  No .env file found or failed to read it: %v", err)
	}

	conf := &Config{
		DB: DBConfig{
			Host:     viper.GetString("DB_HOST"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			DBName:   viper.GetString("DB_NAME"),
			Port:     viper.GetString("DB_PORT"),
			SSLMode:  viper.GetString("DB_SSLMODE"), // optional
		},
	}

	validateConfig(conf)

	AppConfig = conf

	return conf
}

func validateConfig(cfg *Config) {
	if cfg.DB.Host == "" || cfg.DB.User == "" || cfg.DB.DBName == "" {
		log.Fatal("❌ Missing required DB configuration (DB_HOST, DB_USER, DB_NAME)")
	}

}

func (c *DBConfig) Dsn() string {
	ssl := c.SSLMode
	if ssl == "" {
		ssl = "disable"
	}
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=True",
		c.User, c.Password, c.Host, c.Port, c.DBName,
	)
}
