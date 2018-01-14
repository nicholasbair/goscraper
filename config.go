package goscraper

var cs = Configs{
	"indeed": Config{
		"indeed",
		"https://www.indeed.com",
		"https://www.indeed.com/jobs?&limit=50",
		map[string]string{
			"location": "l",
			"q":        "q",
			"not":      "as_not",
			"from_age": "fromage",
		},
		50,
		"&start=",
		"resultCount",
		"#searchCount",
		3,
		".result",
		".jobtitle",
		".company",
		".summary",
		".turnstilelink",
	},
	"dice": Config{
		"dice",
		"https://www.dice.com",
		"https://www.dice.com/jobs/advancedResult.html?&limit=100",
		map[string]string{
			"location": "for_loc", // must have comma after city
			"q":        "for_one",
			"from_age": "postedDate",
			"job_type": "jtype", // Full Time, Part Time, Contracts
			"sort":     "sort",  // relevance, date
			"radius":   "radius",
		},
		100,
		"-startPage-",
		"pageNumber",
		"#posiCountId",
		0,
		".complete-serp-result-div",
		".dice-btn-link",
		".compName",
		".shortdesc",
		".dice-btn-link",
	},
}
