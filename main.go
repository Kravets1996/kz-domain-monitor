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
	"sort"
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

	var domains []api.Domain
	hasError := false

	for i, domainName := range cfg.DomainList {
		domain := api.GetDomainInfo(domainName)

		log.Println(domain.GetMessage())

		hasError = hasError || !domain.IsOk()

		if domain.ShouldSend() {
			domains = append(domains, domain)
		}

		if i < len(cfg.DomainList)-1 {
			// TODO Настраиваемый интервал через .env
			time.Sleep(time.Second * 3)
		}
	}

	sortDomains(domains, cfg.SortOrder)

	var messages []string
	for _, domain := range domains {
		messages = append(messages, domain.GetMessage())
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

func sortDomains(domains []api.Domain, sortOrder string) {
	switch sortOrder {
	case "expiration":
		sort.SliceStable(domains, func(i, j int) bool {
			if domains[i].ExpirationDate == nil {
				return true
			}
			if domains[j].ExpirationDate == nil {
				return false
			}
			return domains[i].ExpirationDate.Before(*domains[j].ExpirationDate)
		})
	case "alphabet":
		sort.SliceStable(domains, func(i, j int) bool {
			return domains[i].Name < domains[j].Name
		})
	}
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
