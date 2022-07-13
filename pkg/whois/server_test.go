package whois

import (
	"DomainMan/pkg/database"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestGetServersWithRfc1036Whois(t *testing.T) {
	suffixes, err := GetServersWithRfc1036Whois()
	if err != nil {
		t.Error(err)
	}
	if len(suffixes) == 0 {
		t.Error("No suffixes returned")
	}
	for _, suffix := range suffixes {
		if suffix.Name == "" {
			t.Error("Suffix name is empty")
		}
		if suffix.WhoisServer == "" {
			t.Error("Suffix whois server is empty")
		}
	}
	type output struct {
		Name        string `json:"name"`
		Mode        uint   `json:"mode"`
		WhoisServer string `json:"whoisServer"`
	}
	var outputs []output
	for _, suffix := range suffixes {
		outputs = append(outputs, output{
			Name:        suffix.Name,
			Mode:        suffix.Mode,
			WhoisServer: suffix.WhoisServer,
		})
	}
	j, _ := json.Marshal(outputs)
	fmt.Println(string(j))
}

func TestUpdateServersWithRfc1036Whois(t *testing.T) {
	os.Setenv("DOMAINMAN_DATABASE_DIALECT", "sqlite")
	os.Setenv("DOMAINMAN_DATABASE_PARAMETER", ":memory:")
	err := database.Connect(true)
	if err != nil {
		t.Error(err)
	}
	err = UpdateServersWithRfc1036Whois(true)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateServersWithPreset(t *testing.T) {
	os.Setenv("DOMAINMAN_DATABASE_DIALECT", "sqlite")
	os.Setenv("DOMAINMAN_DATABASE_PARAMETER", ":memory:")
	err := database.Connect(true)
	if err != nil {
		t.Error(err)
	}
	err = UpdateServersWithPreset(true)
	if err != nil {
		t.Error(err)
	}
}
