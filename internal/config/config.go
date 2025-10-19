package config

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DB      DBConfig
	JWT     JwtConfig
	Xendit  XenditConfig
	AppMain AppMain
}

type AppMain struct {
	PORT               string
	BaseUrl            string
	Domain             string
	GoggleClientId     string
	GoggleClientSecret string
	GoggleRedirectUrl  string
}

type DBConfig struct {
	Host     string
	Password string
	User     string
	DBName   string
	Port     string
	SSLMode  string
}

type JwtConfig struct {
	SecretKey      string
	AccessTokenExp time.Duration
}

type XenditConfig struct {
	APIKey     string
	WebhookKey string
	BaseURL    string
}

var AppConfig *Config

func LoadConfig() *Config {
	viper.SetConfigFile("config.yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/app/")

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("⚠️  No .env file found or failed to read it: %v", err)
	}

	conf := &Config{
		AppMain: AppMain{
			PORT:               viper.GetString("app.port"),
			BaseUrl:            viper.GetString("app.base_url"),
			Domain:             viper.GetString("app.domain"),
			GoggleClientId:     viper.GetString("app.google_client_id"),
			GoggleClientSecret: viper.GetString("app.google_client_secret"),
			GoggleRedirectUrl:  viper.GetString("app.google_redirect_url"),
		},
		DB: DBConfig{
			Host:     viper.GetString("database.host"),
			User:     viper.GetString("database.user"),
			Password: viper.GetString("database.password"),
			DBName:   viper.GetString("database.dbname"),
			Port:     viper.GetString("database.port"),
			SSLMode:  viper.GetString("database.sslmode"),
		},
		JWT: JwtConfig{
			SecretKey:      viper.GetString("jwt.secret_key"),
			AccessTokenExp: viper.GetDuration("jwt.access_token_exp"),
		},
		Xendit: XenditConfig{
			APIKey:     viper.GetString("xendit.api_key"),
			WebhookKey: viper.GetString("xendit.webhook_key"),
			BaseURL:    viper.GetString("xendit.base_url"),
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

	if cfg.JWT.SecretKey == "" {
		log.Fatal("❌ Missing required JWT configuration (SecretKey)")
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

func (c *JwtConfig) GetSecretKey() string {
	return c.SecretKey
}

func (c *JwtConfig) GetExpTimeAccessToken() time.Duration {
	return c.AccessTokenExp
}
func (c *XenditConfig) GetBaseUrl() string {
	return c.BaseURL
}

func (c *XenditConfig) GetApiKey() string {
	return c.APIKey
}

func (c *XenditConfig) GetWebhookKey() string {
	return c.WebhookKey
}

func (c *AppMain) GetDomain() string {
	return c.Domain
}

func (c *AppMain) GetPort() string {
	return c.PORT
}
func (c *AppMain) GetBaseUrl() string {
	return c.BaseUrl
}
