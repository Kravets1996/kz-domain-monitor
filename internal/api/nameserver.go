package api

import "strings"

// kzSecondLevelZones - публичные зоны второго уровня под .kz, в которых
// регистрируемый домен находится на третьем уровне (например example.com.kz).
var kzSecondLevelZones = map[string]bool{
	"com.kz": true,
	"org.kz": true,
	"net.kz": true,
	"edu.kz": true,
	"gov.kz": true,
	"mil.kz": true,
}

// registrableDomain сводит хост NS-сервера к домену, который нужно продлевать,
// чтобы NS продолжал работать:
//
//	"ns1.hoster.kz"      -> "hoster.kz"
//	"dns.example.com.kz" -> "example.com.kz"
//
// Сам по себе домен (без поддоменов) возвращается как есть.
func registrableDomain(host string) string {
	host = strings.Trim(strings.ToLower(strings.TrimSpace(host)), ".")
	if host == "" {
		return ""
	}

	labels := strings.Split(host, ".")
	if len(labels) < 2 {
		return host
	}

	// example.com.kz: регистрируемый домен на третьем уровне.
	if len(labels) >= 3 {
		lastTwo := strings.Join(labels[len(labels)-2:], ".")
		if kzSecondLevelZones[lastTwo] {
			return strings.Join(labels[len(labels)-3:], ".")
		}
	}

	return strings.Join(labels[len(labels)-2:], ".")
}
