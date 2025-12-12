package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"kz-domain-monitor/internal/config"
	"log"
	"net/http"
	"time"
)

var client = http.Client{
	Timeout: time.Second * 10,
}

func GetDomainInfo(domainName string) (string, bool) {
	cfg := config.GetConfig()

	query := fmt.Sprintf(`query {
		domains {
			whois {
				whois(domain:"%s") {
					available 
					info {
						domain {
							exDate
						}
					}
				}
			}
		}
	}`, domainName)

	response, err := sendRequest("https://console.ps.kz/domains/graphql", GraphQLRequest{Query: query})

	if err != nil {
		return "❗️ " + err.Error(), false
	}

	exDate, err := time.Parse(time.RFC3339, response.Data.Domains.Whois.Whois.Info.Domain.ExDate)
	if err != nil {
		return "❗️ Error parsing date: " + err.Error(), false
	}

	diff := exDate.Sub(time.Now())
	days := int64(diff.Hours() / 24)

	var icon string
	var result bool

	if days > cfg.DaysToExpire {
		icon = "✅"
		result = true
	} else {
		icon = "❗️"
		result = false
	}

	message := fmt.Sprintf("%s %d дней - %s", icon, days, domainName)

	log.Println(message)

	return message, result
}

func sendRequest(url string, query GraphQLRequest) (*GraphQLResponse, error) {
	var (
		response    *http.Response
		err         error
		gqlResponse GraphQLResponse
	)

	jsonBody, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-User-Token", config.GetConfig().PSApiToken)

	response, err = retry(request)
	if err != nil {
		log.Printf("HTTP request error: %s", err)
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Printf("Request status: %s", err)
		return nil, fmt.Errorf("request status error: %s", response.Status)
	}

	err = json.NewDecoder(response.Body).Decode(&gqlResponse)
	if err != nil {
		log.Println("Failed to parse JSON response:", err.Error())
		return nil, err
	}

	return &gqlResponse, nil
}

func retry(r *http.Request) (*http.Response, error) {
	var (
		response      *http.Response
		err           error
		retries       = 3
		retryInterval = time.Second * 10
	)

	for retries > 0 {
		response, err = client.Do(r)

		if err == nil {
			break
		}

		retries--

		if retries > 0 {
			time.Sleep(retryInterval)
		}
	}

	return response, err
}
