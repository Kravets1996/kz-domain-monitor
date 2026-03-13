package config

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var Configuration Config

type Config struct {
	PSApiToken     string
	DomainProvider string
	DomainList     []string
	DaysToExpire   int64
	SendSuccess    bool
	SendOnlyErrors bool
	RequestDelay   time.Duration
	SortOrder      string
	Telegram       TelegramConfig
	Slack          SlackConfig
	Email          EmailConfig
	Webhook        WebhookConfig
}

// jsonDomainEntry represents either a domain entry or a group in the JSON config.
type jsonDomainEntry struct {
	Domain string            `json:"domain"`
	Title  string            `json:"title"`
	Items  []jsonDomainEntry `json:"items"`
}

// loadDomainsFromJSON reads a JSON config file and extracts the flat list of domain names.
func loadDomainsFromJSON(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var entries []jsonDomainEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return extractDomains(entries), nil
}

func extractDomains(entries []jsonDomainEntry) []string {
	var domains []string
	for _, e := range entries {
		if e.Domain != "" {
			domains = append(domains, strings.TrimSpace(e.Domain))
		}
		if len(e.Items) > 0 {
			domains = append(domains, extractDomains(e.Items)...)
		}
	}
	return domains
}

type TelegramConfig struct {
	Enabled  bool
	BotToken string
	ChatID   string
}

type SlackConfig struct {
	Enabled    bool
	WebhookURL string
}

type EmailConfig struct {
	Enabled  bool
	Host     string
	Port     string
	Username string
	Password string
	From     string
	To       []string
}

type WebhookConfig struct {
	Enabled bool
	URL     string
}

func loadDomains() []string {
	jsonFile := os.Getenv(`DOMAIN_CONFIG_FILE`)
	envList := os.Getenv(`DOMAIN_LIST`)

	if jsonFile != "" {
		if envList != "" {
			log.Println("WARNING: Both DOMAIN_CONFIG_FILE and DOMAIN_LIST are set. DOMAIN_CONFIG_FILE takes priority.")
		}
		domains, err := loadDomainsFromJSON(jsonFile)
		if err != nil {
			panic("Failed to load DOMAIN_CONFIG_FILE: " + err.Error())
		}
		return domains
	}

	if envList == "" {
		panic("Environment variable DOMAIN_LIST is not set")
	}
	return strings.Split(envList, ",")
}

func Init() {
	daysToExpireInt, _ := strconv.ParseInt(getEnv(`DAYS_TO_EXPIRE`, "5"), 10, 64)
	requestDelayInt, _ := strconv.ParseInt(getEnv(`REQUEST_DELAY`, "3"), 10, 64)
	domainProvider := getEnv(`DOMAIN_PROVIDER`, "rdap")

	psApiToken := ""
	if domainProvider == "pskz" {
		psApiToken = getEnvStrict(`PS_GRAPHQL_TOKEN`)
	} else {
		psApiToken = os.Getenv(`PS_GRAPHQL_TOKEN`)
	}

	domainList := loadDomains()

	Configuration = Config{
		PSApiToken:     psApiToken,
		DomainProvider: domainProvider,
		DomainList:     domainList,
		DaysToExpire:   daysToExpireInt,
		SendSuccess:    getEnv(`SEND_ON_SUCCESS`, "true") == "true",
		SendOnlyErrors: getEnv(`SEND_ONLY_ERRORS`, "false") == "true",
		SortOrder:      getEnv(`SORT_ORDER`, "default"),
		RequestDelay:   time.Second * time.Duration(requestDelayInt),
		Telegram: TelegramConfig{
			Enabled:  getEnv(`TELEGRAM_ENABLED`, "true") == "true",
			BotToken: os.Getenv(`TELEGRAM_BOT_TOKEN`),
			ChatID:   os.Getenv(`TELEGRAM_CHAT_ID`),
		},
		Slack: SlackConfig{
			Enabled:    getEnv(`SLACK_ENABLED`, "false") == "true",
			WebhookURL: os.Getenv(`SLACK_WEBHOOK_URL`),
		},
		Email: EmailConfig{
			Enabled:  getEnv(`EMAIL_ENABLED`, "false") == "true",
			Host:     os.Getenv(`EMAIL_HOST`),
			Port:     getEnv(`EMAIL_PORT`, "465"),
			Username: os.Getenv(`EMAIL_USERNAME`),
			Password: os.Getenv(`EMAIL_PASSWORD`),
			From:     os.Getenv(`EMAIL_FROM`),
			To:       splitAndTrim(os.Getenv(`EMAIL_TO`)),
		},
		Webhook: WebhookConfig{
			Enabled: getEnv(`WEBHOOK_ENABLED`, "false") == "true",
			URL:     os.Getenv(`WEBHOOK_URL`),
		},
	}

	if Configuration.Telegram.Enabled {
		if Configuration.Telegram.BotToken == "" || Configuration.Telegram.ChatID == "" {
			panic("Telegram config is not set")
		}
	}

	if Configuration.Slack.Enabled {
		if Configuration.Slack.WebhookURL == "" {
			panic("Slack webhook URL is not set")
		}
	}

	if Configuration.Email.Enabled {
		if Configuration.Email.Host == "" || Configuration.Email.Username == "" || Configuration.Email.From == "" || len(Configuration.Email.To) == 0 {
			panic("Email config is not set")
		}
	}

	if Configuration.Webhook.Enabled {
		if Configuration.Webhook.URL == "" {
			panic("Webhook URL is not set")
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

func splitAndTrim(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
