package main

import (
	"fmt"
	"kz-domain-monitor/internal/api"
	"kz-domain-monitor/internal/config"
	"kz-domain-monitor/internal/notification"
	"kz-domain-monitor/internal/storage"
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

	store := openHistoryStore(cfg.HistoryDBPath)
	if store != nil {
		defer store.Close()
	}

	var domains []api.Domain
	hasError := false

	for _, domainName := range cfg.DomainList {
		domain := checker.Check(domainName)

		enrichWithHistory(store, &domain)

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

// openHistoryStore открывает базу данных истории сроков истечения рядом с бинарником.
// История - вспомогательная функция, поэтому при ошибке открытия мониторинг
// продолжает работать без отметок о продлении.
func openHistoryStore(path string) *storage.Store {
	if path == "" {
		path = storage.DefaultPath()
	}

	store, err := storage.Open(path)
	if err != nil {
		log.Printf("История: не удалось открыть базу данных: %v", err)
		return nil
	}
	return store
}

// enrichWithHistory подставляет домену (и его NS-доменам) срок истечения из прошлой
// проверки, после чего сохраняет текущий срок в историю.
func enrichWithHistory(store *storage.Store, domain *api.Domain) {
	if store == nil {
		return
	}

	applyHistory(store, domain)
	for i := range domain.NSDomains {
		applyHistory(store, &domain.NSDomains[i])
	}
}

func applyHistory(store *storage.Store, domain *api.Domain) {
	if domain.ExpirationDate == nil {
		return
	}

	prev, err := store.GetExpiration(domain.Name)
	if err != nil {
		log.Printf("История: не удалось прочитать %s: %v", domain.Name, err)
		return
	}
	domain.PreviousExpirationDate = prev

	if err := store.SetExpiration(domain.Name, *domain.ExpirationDate); err != nil {
		log.Printf("История: не удалось сохранить %s: %v", domain.Name, err)
	}
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
