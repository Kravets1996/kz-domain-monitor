package api

import (
	"encoding/json"
	"testing"
)

func TestParse(t *testing.T) {
	var response GraphQLResponse
	jsonString := `{"data":{"domains":{"whois":{"whois":{"available":true,"info":{"domain":{"exDate":"2006-01-02T15:04:05Z07:00"}}}}}}}`

	err := json.Unmarshal([]byte(jsonString), &response)

	if err != nil {
		t.Fatal(err)
	}

	if !response.IsAvailable() {
		t.Error("IsAvailable should be true")
	}

	if response.GetExpirationDate() != "2006-01-02T15:04:05Z07:00" {
		t.Error("GetExpirationDate should be 2006-01-02T15:04:05Z07:00")
	}
}
