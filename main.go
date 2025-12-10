package main

import (
	"kz-domain-monitor/internal/api"
	"kz-domain-monitor/internal/config"
	"kz-domain-monitor/internal/notification"
	"log"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Println(`DotEnv file not found, using OS environment variables.`)
	}

	config.Init()
	cfg := config.GetConfig()

	var messages []string
	var hasError bool

	for _, domain := range cfg.DomainList {
		message, result := api.GetDomainInfo(domain)

		hasError = hasError || !result

		messages = append(messages, message)

		// TODO Настраиваемый интервал через .env
		time.Sleep(time.Second * 3)
	}

	if !hasError && !cfg.SendSuccess {
		return
	}

	notification.SendNotification(messages, hasError)

	// TODO Exit code
}
