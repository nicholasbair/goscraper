package goscraper

type Config struct {
	Provider              string
	BaseURL               string
	TemplateURL           string
	ResultsPerPage        int
	PaginationURL         string
	PaginationType        string // resultCount or pageNumber
	SelectorResultsNumber string
	ResultsNumberIndex    int
	SelectorResultDiv     string
	SelectorTitle         string
	SelectorCompany       string
	SelectorDesc          string
	SelectorURL           string
}

// type Configs []Config
type Configs map[string]Config

type Job struct {
	Title       string
	Company     string
	Description string
	URL         string
	Provider    string
}

type Jobs []Job

type requestURL struct {
	location string
}
