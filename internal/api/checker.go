package api

import (
	"strings"
	"time"
)

// Checker проверяет домены вместе с .kz-доменами, которым принадлежат их NS-серверы.
// Результаты кэшируются, поэтому один и тот же домен не запрашивается дважды,
// а между реальными сетевыми запросами выдерживается задержка во избежание
// ограничений со стороны API.
type Checker struct {
	cache     map[string]Domain
	delay     time.Duration
	requested bool
	// fetch получает информацию о домене (по умолчанию GetDomainInfo);
	// выделено в поле, чтобы логику кэширования можно было тестировать без сети.
	fetch func(string) Domain
}

// NewChecker создаёт Checker, выдерживающий delay между сетевыми запросами.
func NewChecker(delay time.Duration) *Checker {
	return &Checker{
		cache: make(map[string]Domain),
		delay: delay,
		fetch: GetDomainInfo,
	}
}

// Check возвращает информацию о домене, дополненную информацией о .kz-доменах,
// которым принадлежат его NS-серверы.
func (c *Checker) Check(domainName string) Domain {
	domain := c.lookup(domainName)

	// Сам домен помечаем как обработанный, чтобы не мониторить рекурсивные
	// (in-bailiwick) NS вида ns1.example.kz для example.kz.
	seen := map[string]bool{strings.ToLower(domainName): true}

	for _, ns := range domain.Nameservers {
		nsDomain := registrableDomain(ns)

		if nsDomain == "" || seen[nsDomain] {
			continue
		}
		seen[nsDomain] = true

		// Мониторить мы умеем только .kz-домены.
		if !strings.HasSuffix(nsDomain, ".kz") {
			continue
		}

		domain.NSDomains = append(domain.NSDomains, c.lookup(nsDomain))
	}

	return domain
}

// lookup возвращает домен из кэша либо запрашивает его, применяя задержку
// только перед реальным сетевым запросом (но не перед первым).
func (c *Checker) lookup(domainName string) Domain {
	if d, ok := c.cache[domainName]; ok {
		return d
	}

	if c.requested {
		time.Sleep(c.delay)
	}
	c.requested = true

	d := c.fetch(domainName)
	c.cache[domainName] = d

	return d
}
