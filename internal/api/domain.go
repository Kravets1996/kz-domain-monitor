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
	// NSDomains - домены, которым принадлежат NS-серверы домена.
	// Их непродление также ломает работу сайта, поэтому они мониторятся отдельно.
	NSDomains []Domain
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

	if domain.isExpired() {
		return "❗️"
	}

	if domain.isCloseToExpire() {
		return "⚠️"
	}

	return "✅"
}

func (domain Domain) IsOk() bool {
	if domain.Error != nil || domain.ExpirationDate == nil {
		return false
	}

	return !domain.IsAvailable && !domain.isCloseToExpire()
}

// IsHealthy возвращает true, только если в порядке и сам домен, и все домены его NS-серверов.
func (domain Domain) IsHealthy() bool {
	if !domain.IsOk() {
		return false
	}

	for _, ns := range domain.NSDomains {
		if !ns.IsOk() {
			return false
		}
	}

	return true
}

func (domain Domain) ShouldSend() bool {
	// Ошибки (в том числе по NS-доменам) отправляются всегда.
	if !domain.IsHealthy() {
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

// GetMessages возвращает сообщение домена, а следом - сообщения по доменам его NS-серверов
// с отступом, формируя вложенный вид:
//
//	example.kz
//	  ↳ ns-домен1
//	  ↳ ns-домен2
func (domain Domain) GetMessages() []string {
	messages := []string{domain.GetMessage()}

	for _, ns := range domain.NSDomains {
		messages = append(messages, "  ↳ "+ns.GetMessage())
	}

	return messages
}
