package api

import "testing"

func newTestChecker(fetch func(string) Domain) *Checker {
	c := NewChecker(0)
	c.fetch = fetch
	return c
}

// Один и тот же NS-домен, встречающийся у нескольких NS-серверов, запрашивается один раз.
func TestChecker_DeduplicatesNameservers(t *testing.T) {
	calls := map[string]int{}
	checker := newTestChecker(func(name string) Domain {
		calls[name]++
		if name == "example.kz" {
			return Domain{Name: name, Nameservers: []string{"ns1.hoster.kz", "ns2.hoster.kz", "ns3.hoster.kz"}}
		}
		return Domain{Name: name}
	})

	domain := checker.Check("example.kz")

	if len(domain.NSDomains) != 1 {
		t.Fatalf("expected 1 NS domain, got %d: %+v", len(domain.NSDomains), domain.NSDomains)
	}
	if domain.NSDomains[0].Name != "hoster.kz" {
		t.Errorf("expected NS domain hoster.kz, got %s", domain.NSDomains[0].Name)
	}
	if calls["hoster.kz"] != 1 {
		t.Errorf("hoster.kz should be fetched once, got %d", calls["hoster.kz"])
	}
}

// Рекурсивные (in-bailiwick) NS вида ns1.example.kz не мониторятся отдельно.
func TestChecker_SkipsRecursiveNameservers(t *testing.T) {
	checker := newTestChecker(func(name string) Domain {
		return Domain{Name: name, Nameservers: []string{"ns1.example.kz", "ns2.example.kz"}}
	})

	domain := checker.Check("example.kz")

	if len(domain.NSDomains) != 0 {
		t.Fatalf("recursive NS should not produce NS domains, got %+v", domain.NSDomains)
	}
}

// Не-.kz NS-серверы пропускаются, т.к. их не получится проверить через .kz-провайдер.
func TestChecker_SkipsNonKzNameservers(t *testing.T) {
	checker := newTestChecker(func(name string) Domain {
		return Domain{Name: name, Nameservers: []string{"ns1.cloudflare.com", "ns2.cloudflare.com"}}
	})

	domain := checker.Check("example.kz")

	if len(domain.NSDomains) != 0 {
		t.Fatalf("non-kz NS should be skipped, got %+v", domain.NSDomains)
	}
}

// Кэш переиспользуется между разными проверяемыми доменами.
func TestChecker_CachesAcrossDomains(t *testing.T) {
	calls := map[string]int{}
	checker := newTestChecker(func(name string) Domain {
		calls[name]++
		switch name {
		case "a.kz", "b.kz":
			return Domain{Name: name, Nameservers: []string{"ns1.hoster.kz"}}
		default:
			return Domain{Name: name}
		}
	})

	checker.Check("a.kz")
	checker.Check("b.kz")

	if calls["hoster.kz"] != 1 {
		t.Errorf("shared NS domain hoster.kz should be fetched once across domains, got %d", calls["hoster.kz"])
	}
}
