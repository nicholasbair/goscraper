package goscraper

import (
	"fmt"
	"testing"
)

// func TestScrapeDice(t *testing.T) {
// 	p := map[string][]string{
// 		"location": []string{"denver,", "co"},
// 		"q":        []string{"customer", "success", "manager"},
// 		"from_age": []string{"7"},
// 		"job_type": []string{"Full", "Time"},
// 		"sort":     []string{"relevance"},
// 		"radius":   []string{"30"},
// 		"provider": []string{"dice"},
// 	}
// 	r := Scrape(p)
// 	if len(r) == 0 {
// 		t.Error("Expected results of scrape to be greater than zero")
// 	}
// }

func TestScrapeIndeed(t *testing.T) {
	p := map[string][]string{
		"location": []string{"denver", "co"},
		"q":        []string{"customer", "success", "manager"},
		"from_age": []string{"7"},
		"provider": []string{"indeed"},
	}
	r := Scrape(p)
	fmt.Println(len(r))
	if len(r) == 0 {
		t.Error("Expected results of scrape to be greater than zero")
	}
}

func TestGetConfigs(t *testing.T) {
	r := GetConfigs()
	if len(r) == 0 {
		t.Error("Expected results of GetConfigs to be greater than zero")
	}
}

func TestGetConfig(t *testing.T) {
	r := GetConfig("indeed")
	if r.Provider != "indeed" {
		t.Error("Expected results of GetConfig to be indeed config")
	}
}
