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
	hasError := false

	for i, domainName := range cfg.DomainList {
		domain := api.GetDomainInfo(domainName)

		message := domain.GetMessage()

		log.Println(message)

		hasError = hasError || !domain.IsOk()

		if domain.ShouldSend() {
			messages = append(messages, message)
		}

		if i < len(cfg.DomainList)-1 {
			time.Sleep(cfg.RequestDelay)
		}
	}

	if !hasError && !cfg.SendSuccess {
		return
	}

	notification.SendNotification(messages, hasError)

	if hasError {
		os.Exit(1)
	}

	os.Exit(0)
}
