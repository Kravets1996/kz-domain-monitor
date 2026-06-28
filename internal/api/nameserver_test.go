package api

import "testing"

func TestRegistrableDomain(t *testing.T) {
	cases := map[string]string{
		"ns1.hoster.kz":      "hoster.kz",
		"NS2.HOSTER.KZ":      "hoster.kz",
		"ns1.hoster.kz.":     "hoster.kz",
		"  ns1.hoster.kz  ":  "hoster.kz",
		"dns.example.com.kz": "example.com.kz",
		"hoster.kz":          "hoster.kz",
		"a.b.c.nitec.kz":     "nitec.kz",
		"ns1.cloudflare.com": "cloudflare.com",
		"localhost":          "localhost",
		"":                   "",
	}

	for input, expected := range cases {
		if got := registrableDomain(input); got != expected {
			t.Errorf("registrableDomain(%q) = %q, want %q", input, got, expected)
		}
	}
}
