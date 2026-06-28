package api

import (
	"kz-domain-monitor/internal/config"
	"os"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	config.Configuration = config.Config{
		DaysToExpire: 15,
	}
	os.Exit(m.Run())
}

func TestDomain_IsOk_OK(t *testing.T) {
	domain := getBasicDomain()

	if !domain.IsOk() {
		t.Fatal("domain SHOULD be ok", domain)
	}
}

func TestDomain_IsOk_Available(t *testing.T) {
	domain := getBasicDomain()
	domain.IsAvailable = true

	if domain.IsOk() {
		t.Fatal("domain SHOULD NOT be ok", domain)
	}
}

func TestDomain_IsOk_CloseToExpire(t *testing.T) {
	domain := getBasicDomain()

	domain.ExpirationDate = days(10)

	if domain.IsOk() {
		t.Fatal("domain SHOULD NOT be ok", domain)
	}
}

func TestDomain_IsOk_Expired(t *testing.T) {
	domain := getBasicDomain()

	domain.ExpirationDate = days(-10)

	if domain.IsOk() {
		t.Fatal("domain SHOULD NOT be ok", domain)
	}
}

func TestDomain_IsOk_NoDate(t *testing.T) {
	domain := getBasicDomain()

	domain.ExpirationDate = nil

	if domain.IsOk() {
		t.Fatal("domain SHOULD NOT be ok", domain)
	}
}

func TestDomain_GetMessage_OK(t *testing.T) {
	domain := getBasicDomain()

	message := domain.GetMessage()
	exampleMessage := "✅ 90 дней - example.kz"

	if message != exampleMessage {
		t.Fatal("wrong message", message, exampleMessage)
	}
}

func TestDomain_GetMessage_Available(t *testing.T) {
	domain := getBasicDomain()
	domain.IsAvailable = true

	message := domain.GetMessage()
	exampleMessage := "❌ Домен доступен для регистрации: example.kz"

	if message != exampleMessage {
		t.Fatal("wrong message", message, exampleMessage)
	}
}

func TestDomain_GetMessage_CloseToExpire(t *testing.T) {
	domain := getBasicDomain()

	domain.ExpirationDate = days(10)

	message := domain.GetMessage()
	exampleMessage := "⚠️ 10 дней - example.kz"

	if message != exampleMessage {
		t.Fatal("wrong message", message, exampleMessage)
	}
}

func TestDomain_GetMessage_Expired(t *testing.T) {
	domain := getBasicDomain()

	domain.ExpirationDate = days(-10)

	message := domain.GetMessage()
	exampleMessage := "❗️ -10 дней - example.kz"

	if message != exampleMessage {
		t.Fatal("wrong message", message, exampleMessage)
	}
}

func TestDomain_GetMessage_NoDate(t *testing.T) {
	domain := getBasicDomain()

	domain.ExpirationDate = nil

	message := domain.GetMessage()
	exampleMessage := "❗️ Дата истечения оплаты домена example.kz недоступна"

	if message != exampleMessage {
		t.Fatal("wrong message", message, exampleMessage)
	}
}

func TestDomain_IsRenewed(t *testing.T) {
	older := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	newer := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	cases := []struct {
		name string
		prev *time.Time
		cur  *time.Time
		want bool
	}{
		{"renewed", &older, &newer, true},
		{"unchanged", &older, &older, false},
		{"no previous", nil, &newer, false},
		{"no current", &older, nil, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			domain := Domain{ExpirationDate: c.cur, PreviousExpirationDate: c.prev}
			if got := domain.IsRenewed(); got != c.want {
				t.Fatalf("IsRenewed() = %v, want %v", got, c.want)
			}
		})
	}
}

func TestDomain_GetMessage_Renewed(t *testing.T) {
	prev := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	domain := getBasicDomain()
	domain.PreviousExpirationDate = &prev

	message := domain.GetMessage()

	if !strings.Contains(message, "🔄 продление") {
		t.Fatalf("expected renewal marker in message, got %q", message)
	}
	expected := "01.01.2025 → " + domain.ExpirationDate.Format(dateLayout)
	if !strings.Contains(message, expected) {
		t.Fatalf("expected %q in message, got %q", expected, message)
	}
}

func TestDomain_GetMessage_NotRenewed_NoMarker(t *testing.T) {
	domain := getBasicDomain()

	if strings.Contains(domain.GetMessage(), "продление") {
		t.Fatalf("did not expect renewal marker, got %q", domain.GetMessage())
	}
}

func getBasicDomain() Domain {
	return Domain{
		Name:           "example.kz",
		IsAvailable:    false,
		ExpirationDate: days(90),
	}
}

func days(n time.Duration) *time.Time {
	offset := time.Hour * 12
	if n < 0 {
		offset = -offset
	}
	t := time.Now().Add(time.Hour*24*n + offset)
	return &t
}
