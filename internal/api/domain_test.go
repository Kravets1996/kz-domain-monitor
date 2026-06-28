package api

import (
	"kz-domain-monitor/internal/config"
	"os"
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
