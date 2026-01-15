package main

import (
	"fmt"
	"kz-domain-monitor/internal/api"
	"kz-domain-monitor/internal/config"
	"kz-domain-monitor/internal/notification"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/fynelabs/selfupdate"
	"github.com/joho/godotenv"
)

var Version = "dev"

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "version" {
			printVersion()
			return
		}

		if os.Args[1] == "update" {
			update()
			printVersion()

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
			// TODO Настраиваемый интервал через .env
			time.Sleep(time.Second * 3)
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

func printVersion() {
	fmt.Printf("kz-domain-monitor version %s\n", Version)
}

func update() {
	url := "https://github.com/Kravets1996/kz-domain-monitor/releases/latest/download/kz-domain-monitor"

	if runtime.GOOS == "windows" {
		url += ".exe"
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Update failed: %v", err)

		return
	}
	defer resp.Body.Close()

	err = selfupdate.Apply(resp.Body, selfupdate.Options{})

	if err != nil {
		log.Fatalf("Update failed: %v", err)
		return
	}

	fmt.Println("Update successful")
}
