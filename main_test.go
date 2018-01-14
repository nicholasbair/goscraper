package goscraper

import (
	"testing"
)

func TestScrape(t *testing.T) {
	p := map[string][]string{
		"location": []string{"denver", "co"},
		"provider": []string{"indeed"},
	}
	r := Scrape(p)
	if len(r) == 0 {
		t.Error("Expected results of scrape to be greater than zero")
	}
}
