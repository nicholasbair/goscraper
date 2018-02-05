// Package goscraper contains functions and configurations for scraping
package goscraper

import (
	"sync"
	"time"
)

var (
	wg sync.WaitGroup
	ch = make(chan Jobs, 1)
	ce = make(chan error, 1)
)

// Scrape kicks off a scrape
func Scrape(p map[string][]string) (Jobs, error) {
	var err error
	js := Jobs{}

	// Start channel listener for jobs
	go func() {
		for r := range ch {
			js = append(js, r...)
		}
	}()

	// Start channel listener for errors
	go func() {
		for {
			select {
			case err = <-ce:
				return
			}
		}
	}()

	wg.Add(1)
	go cs[p["provider"][0]].doScraping(p)
	wg.Wait()
	// TEMP: allow the last value from the channel to append to js
	time.Sleep(time.Nanosecond)
	return js, err
}

// GetConfigs returns the config struct for each scrape
func GetConfigs() Configs {
	return cs
}

// GetConfig returns a single config struct
func GetConfig(p string) Config {
	return cs[p]
}
