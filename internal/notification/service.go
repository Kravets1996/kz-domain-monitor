package notification

import (
	"fmt"
	"kz-domain-monitor/internal/config"
	"kz-domain-monitor/internal/notification/channels"
	"strings"
)

func SendNotification(results []string, hasError bool) {
	cfg := config.GetConfig()

	if len(results) == 0 {
		return
	}

	message := "До истечения домена осталось: \n"

	message = message + strings.Join(results, "\n")

	if cfg.Telegram.Enabled {
		err := channels.NewTelegramChannel(cfg.Telegram.BotToken, cfg.Telegram.ChatID).Send(message, !hasError)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}

	if cfg.Slack.Enabled {
		err := channels.NewSlackChannel(cfg.Slack.WebhookURL).Send(message)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}

	if cfg.Email.Enabled {
		err := channels.NewEmailChannel(cfg.Email.Host, cfg.Email.Port, cfg.Email.Username, cfg.Email.Password, cfg.Email.From, cfg.Email.To).Send(message)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}

	if cfg.Webhook.Enabled {
		err := channels.NewWebhookChannel(cfg.Webhook.URL).Send(hasError, message)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}
}
