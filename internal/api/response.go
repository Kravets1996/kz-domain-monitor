package api

type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type GraphQLResponse struct {
	Data struct {
		Domains struct {
			Whois struct {
				Whois struct {
					Available bool `json:"available"`
					Info      struct {
						Domain struct {
							ExDate      string `json:"exDate"`
							Nameservers []struct {
								Name string   `json:"name"`
								IP   []string `json:"ip"`
							}
						} `json:"domain"`
					} `json:"info"`
				} `json:"whois"`
			} `json:"whois"`
		} `json:"domains"`
	} `json:"data"`
}

func (r GraphQLResponse) IsAvailable() bool {
	return r.Data.Domains.Whois.Whois.Available
}

func (r GraphQLResponse) GetExpirationDate() string {
	return r.Data.Domains.Whois.Whois.Info.Domain.ExDate
}

func (r GraphQLResponse) GetNameservers() []string {
	var nameservers []string

	for _, ns := range r.Data.Domains.Whois.Whois.Info.Domain.Nameservers {
		nameservers = append(nameservers, ns.Name)
	}

	return nameservers
}
