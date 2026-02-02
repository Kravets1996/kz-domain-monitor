package api

import (
	"fmt"
	"kz-domain-monitor/internal/config"
	"time"
)

type Domain struct {
	Name           string
	IsAvailable    bool
	ExpirationDate *time.Time
	Error          error
	Nameservers    []string
}

func NewDomain(name string, isAvailable bool, expirationDate string, nameservers []string) Domain {
	var datePointer *time.Time
	date, err := time.Parse(time.RFC3339, expirationDate)

	if err != nil {
		datePointer = nil
	} else {
		datePointer = &date
	}

	return Domain{
		Name:           name,
		IsAvailable:    isAvailable,
		ExpirationDate: datePointer,
		Nameservers:    nameservers,
	}
}

func (domain Domain) GetDaysToExpire() int64 {
	diff := time.Until(*domain.ExpirationDate)
	days := int64(diff.Hours() / 24)

	return days
}

func (domain Domain) isCloseToExpire() bool {
	cfg := config.GetConfig()

	return domain.GetDaysToExpire() <= cfg.DaysToExpire
}

func (domain Domain) isExpired() bool {
	return domain.GetDaysToExpire() < 0
}

func (domain Domain) getIcon() string {
	if domain.IsAvailable {
		return "❌"
	}

	if domain.isCloseToExpire() {
		return "⚠️"
	}

	if domain.isExpired() {
		return "❗️"
	}

	return "✅"
}

func (domain Domain) IsOk() bool {
	if domain.Error != nil || domain.ExpirationDate == nil {
		return false
	}

	return !domain.IsAvailable && !domain.isCloseToExpire()
}

func (domain Domain) ShouldSend() bool {
	// Ошибки отправляются всегда.
	if !domain.IsOk() {
		return true
	}

	// Успешные проверки - только если нет флага OnlyErrors
	return !config.GetConfig().SendOnlyErrors
}

func (domain Domain) GetMessage() string {
	if domain.Error != nil {
		return "❗️ " + domain.Error.Error()
	}

	if domain.ExpirationDate == nil {
		return fmt.Sprintf("❗️ Дата истечения оплаты домена %s недоступна", domain.Name)
	}

	if domain.IsAvailable {
		return "❌ Домен доступен для регистрации: " + domain.Name
	}

	return fmt.Sprintf("%s %d дней - %s", domain.getIcon(), domain.GetDaysToExpire(), domain.Name)
}
