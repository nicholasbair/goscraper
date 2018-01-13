package goscraper

import (
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

var wg sync.WaitGroup
var js = jobs{}
var ch = make(chan jobs)

func main() {
	// start := time.Now()

	wg.Add(1)
	go cs[0].scrape()
	go func() {
		for r := range ch {
			js = append(js, r...)
		}
	}()

	wg.Wait()
	// elapsed := time.Since(start)
	// fmt.Println("elapsed =", elapsed)
	// fmt.Println(len(js))
}

func (c config) scrape() {
	defer wg.Done()
	n := getNumResults(c)
	l := getResultLinks(c, n)
	wg.Add(len(l))

	for i, p := range l {
		go getJobData(p, c, i)
	}
}

func getNumResults(c config) int {
	doc, err := goquery.NewDocument(c.uri)
	checkError(err)

	n := doc.Find(c.selectorResultsNumber)
	s := strings.Split(strings.TrimSpace(n.Text()), " ")
	i, err := strconv.Atoi(s[c.resultsNumberIndex])
	return i
}

func getResultLinks(c config, numOfResults int) []string {
	r := []string{}
	var n int

	// Set the capacity of r dynamically instead of resizing

	switch c.paginationType {
	case "resultCount":
		n = c.resultsPerPage
	case "pageNumber":
		n = 1
		numOfResults = numOfResults / c.resultsPerPage
	}

	for i := 0; i < numOfResults; i += n {
		r = append(r, c.uri+c.paginationURL+strconv.Itoa(i))
	}

	return r
}

func getJobData(l string, c config, i int) {
	defer wg.Done()
	resp, err := http.Get(l)
	checkError(err)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	checkError(err)

	j := make(jobs, 0, c.resultsPerPage)

	doc.Find(c.selectorResultDiv).Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find(c.selectorTitle).Text())
		company := strings.TrimSpace(s.Find(c.selectorCompany).Text())
		desc := strings.TrimSpace(s.Find(c.selectorDesc).Text())
		url := s.Find(c.selectorURL)
		u, _ := url.Attr("href")
		j = append(j, job{title, company, desc, c.baseURL + u, c.provider})
	})

	ch <- j
}
