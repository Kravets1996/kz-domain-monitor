package main

import (
	"fmt"
	"kz-domain-monitor/internal/api"
	"kz-domain-monitor/internal/config"
	"kz-domain-monitor/internal/notification"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var Version = "dev"

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "version" {
			fmt.Printf("kz-domain-monitor version %s\n", Version)
			return
		}
	}

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
