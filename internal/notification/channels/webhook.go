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

type WebhookChannel struct {
	url string
}

type WebhookBody struct {
	Message  string `json:"message"`
	HasError bool   `json:"hasError"`
}

func NewWebhookChannel(url string) *WebhookChannel {
	return &WebhookChannel{url: url}
}

func (w *WebhookChannel) Send(hasError bool, message string) error {
	payload, err := json.Marshal(WebhookBody{
		Message:  message,
		HasError: hasError,
	})

	if err != nil {
		return fmt.Errorf("webhook: marshal failed: %w", err)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("webhook: request creation failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("webhook: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("webhook: unexpected status: %s (failed to read response body: %w)", resp.Status, err)
		}
		return fmt.Errorf("webhook: unexpected status: %s, body: %s", resp.Status, string(body))
	}
	return nil

}
