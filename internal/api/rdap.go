package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// RDAPResponse represents the RDAP API response from nic.kz.
type RDAPResponse struct {
	LdhName string      `json:"ldhName"`
	Status  []string    `json:"status"`
	Events  []RDAPEvent `json:"events"`
}

// RDAPEvent represents a single event in the RDAP response.
type RDAPEvent struct {
	Action string `json:"eventAction"`
	Date   string `json:"eventDate"`
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
	req.Header.Set("Accept", "application/rdap+json")

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

	return NewDomain(domainName, false, rdapResp.GetExpirationDate())
}
