package notification

import (
	"fmt"
	"kz-domain-monitor/internal/config"
	"kz-domain-monitor/internal/notification/channels"
	"strings"
)

func SendNotification(results []string, hasError bool) {
	cfg := config.GetConfig()

	message := "\n" + strings.Join(results, "\n")

	if cfg.Telegram.Enabled && cfg.Telegram.BotToken != "" {
		err := channels.NewTelegramChannel(cfg.Telegram.BotToken, cfg.Telegram.ChatID).Send(message, !hasError)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}
}
