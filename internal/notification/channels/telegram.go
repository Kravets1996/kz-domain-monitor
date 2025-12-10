package channels

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type TelegramChannel struct {
	botToken string
	chatID   string
}

func NewTelegramChannel(botToken string, chatID string) *TelegramChannel {
	return &TelegramChannel{
		botToken: botToken,
		chatID:   chatID,
	}
}

func (t TelegramChannel) Send(message string, silent bool) (err error) {
	apiURL := "https://api.telegram.org/bot" + t.botToken + "/sendMessage"

	data := url.Values{}
	data.Set("chat_id", t.chatID)
	data.Set("text", message)
	if silent {
		data.Set("disable_notification", "true")
	}
	resp, err := http.Post(
		apiURL,
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram api error: %s", resp.Status)
	}

	return nil
}
