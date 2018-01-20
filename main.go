// Package goscraper contains functions and configurations for scraping
package goscraper

import (
	"sync"
	"time"
)

var wg sync.WaitGroup
var ch = make(chan Jobs)

// Scrape kicks off a scrape
func Scrape(p map[string][]string) Jobs {
	js := Jobs{}

	// Start channel listener
	go func() {
		for r := range ch {
			js = append(js, r...)
		}
	}()

	wg.Add(1)
	go cs[p["provider"][0]].doScraping(p)
	wg.Wait()
	// TEMP: allow the last value from the channel to append to js
	time.Sleep(time.Nanosecond)
	return js
}

// GetConfigs returns the config struct for each scrape
func GetConfigs() Configs {
	return cs
}

// GetConfig returns a single config struct
func GetConfig(p string) Config {
	return cs[p]
}
