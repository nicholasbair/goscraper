// Package goscraper contains functions and configurations for scraping
package goscraper

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

// TODO
// Fix indeed url, just returning base url

var wg sync.WaitGroup
var ch = make(chan Jobs)

// Scrape kicks off a full scrape
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
	close(ch)
	return js
}

func (c Config) doScraping(p map[string][]string) {
	defer wg.Done()
	u := buildSearchURL(c, p)
	n := getNumResults(c, u)
	// Pagination: use p["page"] to config the correct page links to return
	l := getResultLinks(c, n, u)
	wg.Add(len(l))

	for _, link := range l {
		go getJobData(link, c)
	}
	wg.Wait()
}

func buildSearchURL(c Config, p map[string][]string) string {
	u, err := url.Parse(c.TemplateURL)
	q := u.Query()
	checkError(err)

	// Better to build a helper function that deletes k/v pairs from p that don't patch
	// the query map in c Config
	delete(p, "provider")
	delete(p, "page")

	for k, v := range p {
		q.Set(c.QueryMap[k], strings.Join(v, " "))
	}
	u.RawQuery = q.Encode()
	return u.String()
}

// TODO handle # of results w/ a comma "1,790"
func getNumResults(c Config, u string) int {
	resp, err := http.Get(u)
	checkError(err)
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	checkError(err)

	resp.Body.Close()

	n := doc.Find(c.SelectorResultsNumber)
	s := strings.Split(strings.TrimSpace(n.Text()), " ")
	i, err := strconv.Atoi(s[c.ResultsNumberIndex])
	return i
}

func getResultLinks(c Config, numOfResults int, u string) []string {
	r := make([]string, 0, 200/c.ResultsPerPage)

	var n int
	var limit int

	switch c.PaginationType {
	case "resultCount": // Indeed
		n = c.ResultsPerPage
		limit = 200
	case "pageNumber": // Dice
		n = 1
		// numOfResults = numOfResults / c.ResultsPerPage
		// numOfResults = 200 / c.ResultsPerPage
		limit = 200 / c.ResultsPerPage
	}

	for i := 0; i < limit; i += n {
		r = append(r, u+c.PaginationURL+strconv.Itoa(i))
	}

	return r
}

func getJobData(l string, c Config) {
	defer wg.Done()
	resp, err := http.Get(l)
	checkError(err)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	checkError(err)
	resp.Body.Close()

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
func GetConfigs() Configs {
	return cs
}

// GetConfig returns a single config struct
func GetConfig(p string) Config {
	return cs[p]
}
