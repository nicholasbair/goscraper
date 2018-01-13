// Package goscraper contains functions and configurations for scraping
package goscraper

import (
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

var wg sync.WaitGroup
var js = Jobs{}
var ch = make(chan Jobs)

// Scrape kicks off a full scrape
func Scrape() Jobs {
	// start := time.Now()

	wg.Add(1)
	go cs[0].doScraping()
	go func() {
		for r := range ch {
			js = append(js, r...)
		}
	}()

	wg.Wait()
	// elapsed := time.Since(start)
	// fmt.Println("elapsed =", elapsed)
	// fmt.Println(len(js))
	return js
}

func (c config) doScraping() {
	defer wg.Done()
	n := getNumResults(c)
	l := getResultLinks(c, n)
	wg.Add(len(l))

	for i, p := range l {
		go getJobData(p, c, i)
	}
}

func getNumResults(c config) int {
	doc, err := goquery.NewDocument(c.Uri)
	checkError(err)

	n := doc.Find(c.SelectorResultsNumber)
	s := strings.Split(strings.TrimSpace(n.Text()), " ")
	i, err := strconv.Atoi(s[c.ResultsNumberIndex])
	return i
}

func getResultLinks(c config, numOfResults int) []string {
	r := []string{}
	var n int

	// Set the capacity of r dynamically instead of resizing

	switch c.PaginationType {
	case "resultCount":
		n = c.ResultsPerPage
	case "pageNumber":
		n = 1
		numOfResults = numOfResults / c.ResultsPerPage
	}

	for i := 0; i < numOfResults; i += n {
		r = append(r, c.Uri+c.PaginationURL+strconv.Itoa(i))
	}

	return r
}

func getJobData(l string, c config, i int) {
	defer wg.Done()
	resp, err := http.Get(l)
	checkError(err)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	checkError(err)

	j := make(Jobs, 0, c.ResultsPerPage)

	doc.Find(c.SelectorResultDiv).Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find(c.SelectorTitle).Text())
		company := strings.TrimSpace(s.Find(c.SelectorCompany).Text())
		desc := strings.TrimSpace(s.Find(c.SelectorDesc).Text())
		url := s.Find(c.SelectorURL)
		u, _ := url.Attr("href")
		j = append(j, Job{title, company, desc, c.BaseURL + u, c.Provider})
	})

	ch <- j
}

// GetConfigs returns the config struct for each scrape
func GetConfigs() configs {
	return cs
}
