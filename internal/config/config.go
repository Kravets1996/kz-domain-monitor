package config

import (
	"os"
	"strconv"
	"strings"
)

var Configuration Config

type Config struct {
	PSApiToken     string
	DomainProvider string
	DomainList     []string
	DaysToExpire   int64
	SendSuccess    bool
	SendOnlyErrors bool
	SortOrder      string
	Telegram       TelegramConfig
}

type TelegramConfig struct {
	Enabled  bool
	BotToken string
	ChatID   string
}

func Init() {
	daysToExpireInt, _ := strconv.ParseInt(getEnv(`DAYS_TO_EXPIRE`, "5"), 10, 64)
	domainProvider := getEnv(`DOMAIN_PROVIDER`, "pskz")

	psApiToken := ""
	if domainProvider == "pskz" {
		psApiToken = getEnvStrict(`PS_GRAPHQL_TOKEN`)
	} else {
		psApiToken = os.Getenv(`PS_GRAPHQL_TOKEN`)
	}

	Configuration = Config{
		PSApiToken:     psApiToken,
		DomainProvider: domainProvider,
		DomainList:     strings.Split(getEnvStrict(`DOMAIN_LIST`), ","),
		DaysToExpire:   daysToExpireInt,
		SendSuccess:    getEnv(`SEND_ON_SUCCESS`, "true") == "true",
		SendOnlyErrors: getEnv(`SEND_ONLY_ERRORS`, "false") == "true",
		SortOrder:      getEnv(`SORT_ORDER`, "default"),
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
