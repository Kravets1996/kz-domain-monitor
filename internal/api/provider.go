package api

import "kz-domain-monitor/internal/config"

// Provider defines the interface for domain info providers.
type Provider interface {
	GetDomainInfo(domainName string) Domain
}

// NewProvider returns the appropriate provider based on configuration.
func NewProvider(cfg config.Config) Provider {
	switch cfg.DomainProvider {
	case "rdap":
		return &RDAPProvider{}
	default:
		return &PsKzProvider{}
	}
}

// GetDomainInfo fetches domain information using the configured provider.
func GetDomainInfo(domainName string) Domain {
	return NewProvider(config.GetConfig()).GetDomainInfo(domainName)
}
