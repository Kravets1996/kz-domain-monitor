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

	checker := api.NewChecker(cfg.RequestDelay)

	var domains []api.Domain
	hasError := false

	for _, domainName := range cfg.DomainList {
		domain := checker.Check(domainName)

		for _, message := range domain.GetMessages() {
			log.Println(message)
		}

		hasError = hasError || !domain.IsHealthy()

		if domain.ShouldSend() {
			domains = append(domains, domain)
		}
	}

	sortDomains(domains, cfg.SortOrder)

	var messages []string
	if cfg.SortOrder == "group" && len(cfg.DomainGroups) > 0 {
		messages = buildGroupedMessages(domains, cfg.DomainGroups)
	} else {
		for _, domain := range domains {
			messages = append(messages, domain.GetMessages()...)
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

func buildGroupedMessages(domains []api.Domain, groups []config.DomainGroup) []string {
	domainMap := make(map[string]api.Domain, len(domains))
	for _, d := range domains {
		domainMap[d.Name] = d
	}

	var messages []string
	for index, group := range groups {
		var groupMessages []string
		for _, name := range group.Domains {
			if d, ok := domainMap[name]; ok {
				groupMessages = append(groupMessages, d.GetMessages()...)
			}
		}
		if len(groupMessages) > 0 {
			if group.Title != "" {
				messages = append(messages, group.Title+":")
			}
			messages = append(messages, groupMessages...)

			if index < len(groups)-1 {
				messages = append(messages, "")
			}
		}
	}
	return messages
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
