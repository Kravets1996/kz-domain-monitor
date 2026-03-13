package channels

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"
)

type EmailChannel struct {
	host     string
	port     string
	username string
	password string
	from     string
	to       []string
}

func NewEmailChannel(host, port, username, password, from string, to []string) *EmailChannel {
	return &EmailChannel{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
		to:       to,
	}
}

func (e *EmailChannel) Send(message string) error {
	addr := e.host + ":" + e.port

	subject := "kz-domain-monitor уведомление"
	body := "To: " + strings.Join(e.to, ",") + "\r\n" +
		"From: " + e.from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" +
		message

	auth := smtp.PlainAuth("", e.username, e.password, e.host)

	tlsConfig := &tls.Config{
		ServerName: e.host,
	}

	dialer := &net.Dialer{Timeout: 10 * time.Second}
	conn, err := tls.DialWithDialer(dialer, "tcp", addr, tlsConfig)

	if err != nil {
		// Fallback: try plain SMTP with STARTTLS
		return e.sendSTARTTLS(addr, auth, body)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, e.host)
	if err != nil {
		return fmt.Errorf("email: failed to create smtp client: %w", err)
	}
	defer client.Close()

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("email: auth failed: %w", err)
	}

	if err = client.Mail(e.from); err != nil {
		return fmt.Errorf("email: MAIL FROM failed: %w", err)
	}

	for _, recipient := range e.to {
		if err = client.Rcpt(recipient); err != nil {
			return fmt.Errorf("email: RCPT TO failed for %s: %w", recipient, err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("email: DATA failed: %w", err)
	}

	_, err = w.Write([]byte(body))
	if err != nil {
		return fmt.Errorf("email: write failed: %w", err)
	}

	if err = w.Close(); err != nil {
		return fmt.Errorf("email: close writer failed: %w", err)
	}

	return client.Quit()
}

func (e *EmailChannel) sendSTARTTLS(addr string, auth smtp.Auth, body string) error {
	tlsConfig := &tls.Config{
		ServerName: e.host,
	}

	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("email: dial failed: %w", err)
	}
	defer client.Close()

	if ok, _ := client.Extension("STARTTLS"); ok {
		if err = client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("email: STARTTLS failed: %w", err)
		}
	}

	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("email: auth failed: %w", err)
	}

	if err = client.Mail(e.from); err != nil {
		return fmt.Errorf("email: MAIL FROM failed: %w", err)
	}

	for _, recipient := range e.to {
		if err = client.Rcpt(recipient); err != nil {
			return fmt.Errorf("email: RCPT TO failed for %s: %w", recipient, err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("email: DATA failed: %w", err)
	}

	_, err = w.Write([]byte(body))
	if err != nil {
		return fmt.Errorf("email: write failed: %w", err)
	}

	if err = w.Close(); err != nil {
		return fmt.Errorf("email: close writer failed: %w", err)
	}

	return client.Quit()
}
