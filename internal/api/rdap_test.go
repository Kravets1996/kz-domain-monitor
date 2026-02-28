package api

import (
	"encoding/json"
	"fmt"
	"kz-domain-monitor/internal/config"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRDAPResponse_GetExpirationDate(t *testing.T) {
	var resp RDAPResponse
	jsonString := `{
		"ldhName": "example.kz",
		"status": ["active"],
		"events": [
			{"eventAction": "registration", "eventDate": "2021-01-01T00:00:00Z"},
			{"eventAction": "expiration", "eventDate": "2025-01-01T00:00:00Z"},
			{"eventAction": "last changed", "eventDate": "2022-06-15T00:00:00Z"}
		]
	}`

	if err := json.Unmarshal([]byte(jsonString), &resp); err != nil {
		t.Fatal(err)
	}

	expected := "2025-01-01T00:00:00Z"
	if resp.GetExpirationDate() != expected {
		t.Errorf("expected %s, got %s", expected, resp.GetExpirationDate())
	}
}

func TestRDAPResponse_GetExpirationDate_Missing(t *testing.T) {
	resp := RDAPResponse{
		Events: []RDAPEvent{
			{Action: "registration", Date: "2021-01-01T00:00:00Z"},
		},
	}

	if resp.GetExpirationDate() != "" {
		t.Errorf("expected empty string, got %s", resp.GetExpirationDate())
	}
}

func TestRDAPProvider_GetDomainInfo_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/rdap+json")
		fmt.Fprint(w, `{
			"ldhName": "example.kz",
			"status": ["active"],
			"events": [
				{"eventAction": "expiration", "eventDate": "2027-01-01T00:00:00Z"}
			]
		}`)
	}))
	defer server.Close()

	provider := &RDAPProvider{BaseURL: server.URL}
	domain := provider.GetDomainInfo("example.kz")

	if domain.Error != nil {
		t.Fatalf("unexpected error: %v", domain.Error)
	}
	if domain.IsAvailable {
		t.Error("domain should not be available")
	}
	if domain.ExpirationDate == nil {
		t.Error("expiration date should not be nil")
	}
}

func TestRDAPProvider_GetDomainInfo_Available(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer server.Close()

	provider := &RDAPProvider{BaseURL: server.URL}
	domain := provider.GetDomainInfo("available.kz")

	if domain.Error != nil {
		t.Fatalf("unexpected error: %v", domain.Error)
	}
	if !domain.IsAvailable {
		t.Error("domain should be available (404 response)")
	}
}

func TestRDAPProvider_GetDomainInfo_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}))
	defer server.Close()

	provider := &RDAPProvider{BaseURL: server.URL}
	domain := provider.GetDomainInfo("error.kz")

	if domain.Error == nil {
		t.Error("expected error for 500 response")
	}
}

func TestRDAPProvider_GetDomainInfo_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/rdap+json")
		fmt.Fprint(w, `invalid json`)
	}))
	defer server.Close()

	provider := &RDAPProvider{BaseURL: server.URL}
	domain := provider.GetDomainInfo("example.kz")

	if domain.Error == nil {
		t.Error("expected error for invalid JSON response")
	}
}

func TestNewProvider_RDAP(t *testing.T) {
	cfg := config.Config{DomainProvider: "rdap"}
	provider := NewProvider(cfg)
	if _, ok := provider.(*RDAPProvider); !ok {
		t.Error("expected RDAPProvider for 'rdap' config")
	}
}

func TestNewProvider_PsKz(t *testing.T) {
	cfg := config.Config{DomainProvider: "pskz"}
	provider := NewProvider(cfg)
	if _, ok := provider.(*PsKzProvider); !ok {
		t.Error("expected PsKzProvider for 'pskz' config")
	}
}

func TestNewProvider_Default(t *testing.T) {
	cfg := config.Config{DomainProvider: ""}
	provider := NewProvider(cfg)
	if _, ok := provider.(*PsKzProvider); !ok {
		t.Error("expected PsKzProvider as default when DomainProvider is empty")
	}
}
