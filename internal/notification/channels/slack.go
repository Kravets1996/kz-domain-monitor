package channels

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
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

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.webhookURL, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("slack: request creation failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("slack: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("slack: unexpected status: %s (failed to read response body: %w)", resp.Status, err)
		}
		return fmt.Errorf("slack: unexpected status: %s, body: %s", resp.Status, string(body))
	}
	return nil
}
