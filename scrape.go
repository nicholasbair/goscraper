package goscraper

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (c Config) doScraping(p map[string][]string) {
	defer wg.Done()
	// Pass in 'test' as page number in url
	// handlePageNum should return error
	// and check error should push that into the channel
	// BUT program will continue and break later
	// when it tries to consume pNum
	pNum, err := handlePageNum(p)
	checkError(err)
	u := buildSearchURL(c, p)
	n := getNumResults(c, u)
	l := getResultLinks(c, n, u, pNum)
	wg.Add(len(l))

	for _, link := range l {
		go getJobData(link, c)
	}
}

func handlePageNum(p map[string][]string) (int, error) {
	if p["page"] != nil {
		return strconv.Atoi(p["page"][0])
	}
	return 1, nil
}

func buildSearchURL(c Config, p map[string][]string) string {
	u, err := url.Parse(c.TemplateURL)
	checkError(err)
	q := u.Query()
	cleanMap(c, p)

	for k, v := range p {
		switch {
		// Dice expects a comma after city
		case c.Provider == "dice" && k == "location" && len(p) > 1:
			q.Set(c.QueryMap[k], strings.Join(v, ", "))
		default:
			q.Set(c.QueryMap[k], strings.Join(v, " "))
		}
	}
	u.RawQuery = q.Encode()
	return u.String()
}

// TODO: Handle zero results here
func getNumResults(c Config, u string) int {
	resp, err := http.Get(u)
	checkError(err)
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	checkError(err)

	n := doc.Find(c.SelectorResultsNumber)
	s := cleanString(n.Text())
	i, err := strconv.Atoi(s[c.ResultsNumberIndex])
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

func getJobData(l string, c Config) {
	defer wg.Done()
	resp, err := http.Get(l)
	checkError(err)
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	checkError(err)

	j := make(Jobs, 0, c.ResultsPerPage)

	doc.Find(c.SelectorResultDiv).Each(func(i int, s *goquery.Selection) {
		title := strip(strings.TrimSpace(s.Find(c.SelectorTitle).Text()))
		company := strip(strings.TrimSpace(s.Find(c.SelectorCompany).Text()))
		desc := strip(strings.TrimSpace(s.Find(c.SelectorDesc).Text()))
		url := s.Find(c.SelectorURL)
		u, _ := url.Attr("href")
		j = append(j, Job{title, company, desc, c.BaseURL + u, c.Provider})
	})

	ch <- j
}
