package goscraper

// Config for each provider
type Config struct {
	Provider              string
	BaseURL               string
	TemplateURL           string
	QueryMap              map[string]string
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

// Configs is a collection of configs
type Configs map[string]Config

// Job is a result the Scraper returns to the server
type Job struct {
	Title       string
	Company     string
	Description string
	URL         string
	Provider    string
}

// Jobs is a collection of results
type Jobs []Job
