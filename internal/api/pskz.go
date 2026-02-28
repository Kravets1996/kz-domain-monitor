package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"kz-domain-monitor/internal/config"
	"log"
	"net/http"
	"os"
	"time"
)

var client = http.Client{
	Timeout: time.Second * 10,
}

// PsKzProvider fetches domain info from the ps.kz GraphQL API.
type PsKzProvider struct{}

func (p *PsKzProvider) GetDomainInfo(domainName string) Domain {
	// TODO Проверка что домен .kz

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
		return Domain{Error: err}
	}

	var datePointer *time.Time
	date, err := time.Parse(time.RFC3339, response.GetExpirationDate())

	if err != nil {
		datePointer = nil
	} else {
		datePointer = &date
	}

	return Domain{
		Name:           domainName,
		IsAvailable:    response.IsAvailable(),
		ExpirationDate: datePointer,
	}
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
		log.Printf("Request status: %d", response.StatusCode)

		writeErrorToFile(fmt.Sprintf("Request status error: %d, Body: %s", response.StatusCode, bodyToString(response.Body)))

		return nil, fmt.Errorf("request status error: %d", response.StatusCode)
	}

	err = json.NewDecoder(response.Body).Decode(&gqlResponse)
	if err != nil {
		log.Println("Failed to parse JSON response:", err.Error())
		return nil, fmt.Errorf("failed to parse JSON response: %s", err.Error())
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

		log.Println("Retrying request:", err.Error())

		retries--

		if retries > 0 {
			time.Sleep(retryInterval)
		}
	}

	return response, err
}

func bodyToString(body io.ReadCloser) string {
	bodyBytes := new(bytes.Buffer)
	if _, err := bodyBytes.ReadFrom(body); err != nil {
		log.Printf("Failed to read response body: %v", err)
		return ""
	}
	return bodyBytes.String()
}

func writeErrorToFile(errorMsg string) {
	f, err := os.OpenFile("error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		return
	}
	defer f.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("[%s] %s\n", timestamp, errorMsg)

	if _, err := f.WriteString(logEntry); err != nil {
		log.Printf("Failed to write to log file: %v", err)
	}
}
