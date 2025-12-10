package config

import (
	"os"
	"strconv"
	"strings"
)

var Configuration Config

type Config struct {
	PSApiToken   string
	DomainList   []string
	DaysToExpire int64
	SendSuccess  bool
	Telegram     TelegramConfig
}

type TelegramConfig struct {
	Enabled  bool
	BotToken string
	ChatID   string
}

func Init() {
	daysToExpireInt, _ := strconv.ParseInt(getEnv(`DAYS_TO_EXPIRE`, "5"), 10, 64)

	Configuration = Config{
		PSApiToken:   getEnvStrict(`PS_GRAPHQL_TOKEN`),
		DomainList:   strings.Split(getEnvStrict(`DOMAIN_LIST`), ","),
		DaysToExpire: daysToExpireInt,
		SendSuccess:  getEnv(`SEND_SUCCESS`, "true") == "true",
		Telegram: TelegramConfig{
			Enabled:  getEnv(`TELEGRAM_ENABLED`, "true") == "true",
			BotToken: os.Getenv(`TELEGRAM_BOT_TOKEN`),
			ChatID:   os.Getenv(`TELEGRAM_CHAT_ID`),
		},
	}

	if Configuration.Telegram.Enabled {
		if Configuration.Telegram.BotToken == "" || Configuration.Telegram.ChatID == "" {
			panic("Telegram config is not set")
		}
	}
}

func GetConfig() Config {
	return Configuration
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvStrict(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	panic("Environment variable " + key + " is not set")
}
