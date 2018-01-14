// Package goscraper contains functions and configurations for scraping
package goscraper

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"text/template"

	"github.com/PuerkitoBio/goquery"
)

// TODO
// Fix indeed url, just returning base url
// What query params will we need to build out search URL?
// keywords
// location - city and state
// ...

var wg sync.WaitGroup
var ch = make(chan Jobs)

// Scrape kicks off a full scrape
func Scrape(p map[string][]string) Jobs {
	js := Jobs{}

	wg.Add(1)
	go cs[p["provider"][0]].doScraping(p)
	go func() {
		for r := range ch {
			js = append(js, r...)
		}
	}()

	wg.Wait()
	return js
}

func (c Config) doScraping(p map[string][]string) {
	defer wg.Done()
	r := buildRequest(c, p)
	u := buildSearchURL(c, r)
	fmt.Println(u)
	n := getNumResults(c, u)
	l := getResultLinks(c, n, u)
	wg.Add(len(l))

	for i, p := range l {
		go getJobData(p, c, i)
	}
}

// map[location:[denver co] provider:[dice]]
// https://www.indeed.com/jobs?q=customer+success+manager&l=
func buildSearchURL(c Config, r RequestURL) string {
	var tpl bytes.Buffer
	// t, err := template.New("test").Parse(c.TemplateURL)
	t, err := template.New("test").Parse("https://www.indeed.com/jobs?q=customer+success+manager&l={{.location}}&as_not=travel&fromage=7&limit=50")
	checkError(err)
	t.Execute(&tpl, r)
	return tpl.String()
}

func buildRequest(c Config, p map[string][]string) RequestURL {
	// delete(p, "provider")
	r := RequestURL{
		Location: strings.Join(p["location"], "+"),
	}
	return r
}

func getNumResults(c Config, u string) int {
	doc, err := goquery.NewDocument(u)
	checkError(err)

	n := doc.Find(c.SelectorResultsNumber)
	s := strings.Split(strings.TrimSpace(n.Text()), " ")
	i, err := strconv.Atoi(s[c.ResultsNumberIndex])
	return i
}

func getResultLinks(c Config, numOfResults int, u string) []string {
	// Set the capacity of r dynamically instead of resizing
	r := []string{}
	var n int

	switch c.PaginationType {
	case "resultCount":
		n = c.ResultsPerPage
	case "pageNumber":
		n = 1
		numOfResults = numOfResults / c.ResultsPerPage
	}

	for i := 0; i < numOfResults; i += n {
		r = append(r, u+c.PaginationURL+strconv.Itoa(i))
	}

	return r
}

func getJobData(l string, c Config, i int) {
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
func GetConfigs() Configs {
	return cs
}
