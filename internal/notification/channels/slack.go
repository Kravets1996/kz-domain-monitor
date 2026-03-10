package channels

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type SlackChannel struct {
	webhookURL string
}

func NewSlackChannel(webhookURL string) *SlackChannel {
	return &SlackChannel{webhookURL: webhookURL}
}

func (s *SlackChannel) Send(message string) error {
	payload, err := json.Marshal(map[string]string{"text": message})
	if err != nil {
		return fmt.Errorf("slack: marshal failed: %w", err)
	}

	resp, err := http.Post(s.webhookURL, "application/json", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("slack: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("slack: unexpected status: %s", resp.Status)
	}

	return nil
}
