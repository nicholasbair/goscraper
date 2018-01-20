// Package goscraper contains functions and configurations for scraping
package goscraper

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

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
	// TEMP: allow the last value from the channel to append to js
	time.Sleep(time.Nanosecond)
	return js
}

func (c Config) doScraping(p map[string][]string) {
	defer wg.Done()
	pNum := handlePageNum(p)
	u := buildSearchURL(c, p)
	fmt.Println("buildSearchURL =", u)
	n := getNumResults(c, u)
	l := getResultLinks(c, n, u, pNum)
	wg.Add(len(l))

	for _, link := range l {
		go getJobData(link, c)
	}
}

func handlePageNum(p map[string][]string) int {
	var pNum int

	if p["page"] != nil {
		pNum, _ = strconv.Atoi(p["page"][0])
	} else {
		pNum = 1
	}
	return pNum
}

func buildSearchURL(c Config, p map[string][]string) string {
	u, err := url.Parse(c.TemplateURL)
	checkError(err)
	q := u.Query()

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
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	checkError(err)

	n := doc.Find(c.SelectorResultsNumber)
	s := strings.Split(strings.TrimSpace(n.Text()), " ")
	i, err := strconv.Atoi(s[c.ResultsNumberIndex])
	fmt.Println("n.Text()=", n.Text())
	fmt.Println("s =", s)
	fmt.Println("i =", i)
	checkError(err)

	return i
}

func getResultLinks(c Config, numOfResults int, u string, pageNum int) []string {
	const paginationNum = 200
	var n int
	var limit int
	var index int

	r := make([]string, 0, paginationNum/c.ResultsPerPage)

	switch c.PaginationType {
	case "resultCount": // Indeed
		n = c.ResultsPerPage
		limit = paginationNum * pageNum
		if pageNum == 1 {
			index = 0
		} else {
			index = (pageNum - 1) * paginationNum
		}
	case "pageNumber": // Dice
		n = 1
		index = (paginationNum / c.ResultsPerPage) + 1
		limit = index + pageNum
	}

	for i := index; i < limit; i += n {
		r = append(r, u+c.PaginationURL+strconv.Itoa(i))
	}

	return r
}

// TODO: build helper function for http
func getJobData(l string, c Config) {
	defer wg.Done()
	resp, err := http.Get(l)
	checkError(err)
	defer resp.Body.Close()

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

	fmt.Println(j)

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
