package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// RDAPResponse represents the RDAP API response from nic.kz.
type RDAPResponse struct {
	LdhName     string           `json:"ldhName"`
	Status      []string         `json:"status"`
	Events      []RDAPEvent      `json:"events"`
	Nameservers []RDAPNameserver `json:"nameservers"`
}

// RDAPEvent represents a single event in the RDAP response.
type RDAPEvent struct {
	Action string `json:"eventAction"`
	Date   string `json:"eventDate"`
}

// RDAPNameserver represents a nameserver entry in the RDAP response.
type RDAPNameserver struct {
	LdhName string `json:"ldhName"`
}

// GetNameservers returns the nameserver host names from the RDAP response.
func (r RDAPResponse) GetNameservers() []string {
	var nameservers []string
	for _, ns := range r.Nameservers {
		if ns.LdhName != "" {
			nameservers = append(nameservers, ns.LdhName)
		}
	}
	return nameservers
}

// GetExpirationDate returns the expiration date string from RDAP events.
func (r RDAPResponse) GetExpirationDate() string {
	for _, e := range r.Events {
		if e.Action == "expiration" {
			return e.Date
		}
	}
	return ""
}

// parseRDAPDate parses dates in RFC3339 format (e.g. "2025-05-08T21:59:03Z").
func parseRDAPDate(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

// RDAPProvider fetches domain info from the nic.kz RDAP endpoint.
type RDAPProvider struct {
	BaseURL string // defaults to "https://rdap.nic.kz"
}

func (p *RDAPProvider) baseURL() string {
	if p.BaseURL != "" {
		return p.BaseURL
	}
	return "https://rdap.nic.kz"
}

func (p *RDAPProvider) GetDomainInfo(domainName string) Domain {
	url := fmt.Sprintf("%s/domain/%s", p.baseURL(), domainName)
	return rdapGetDomainInfoFromURL(url, domainName)
}

func rdapGetDomainInfoFromURL(url, domainName string) Domain {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Domain{Error: err}
	}

	resp, err := retry(req)
	if err != nil {
		log.Printf("RDAP HTTP request error: %s", err)
		return Domain{Error: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return Domain{Name: domainName, IsAvailable: true}
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("RDAP request status: %d", resp.StatusCode)
		writeErrorToFile(fmt.Sprintf("RDAP request status error: %d, Body: %s", resp.StatusCode, bodyToString(resp.Body)))
		return Domain{Error: fmt.Errorf("RDAP request status error: %d", resp.StatusCode)}
	}

	var rdapResp RDAPResponse
	if err := json.NewDecoder(resp.Body).Decode(&rdapResp); err != nil {
		log.Println("Failed to parse RDAP JSON response:", err.Error())
		return Domain{Error: fmt.Errorf("failed to parse RDAP JSON response: %s", err.Error())}
	}

	var datePointer *time.Time

	dateStr := rdapResp.GetExpirationDate()

	date, err := parseRDAPDate(dateStr)

	if err != nil {
		datePointer = nil
	} else {
		datePointer = &date
	}

	return Domain{
		Name:           domainName,
		IsAvailable:    false,
		ExpirationDate: datePointer,
		Nameservers:    rdapResp.GetNameservers(),
	}
}
