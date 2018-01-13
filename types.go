package goscraper

type Config struct {
	Provider              string
	BaseURL               string
	Uri                   string
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
	Url         string
	Provider    string
}

type Jobs []Job
