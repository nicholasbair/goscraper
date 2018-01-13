package goscraper

var cs = configs{
	config{
		"indeed",
		"https://www.indeed.com",
		"https://www.indeed.com/jobs?q=customer+success+manager&l=denver,+co&as_not=travel&fromage=7&limit=50",
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
	config{
		"dice",
		"https://www.dice.com",
		"https://www.dice.com/jobs/advancedResult.html?for_one=customer+success+manager&for_all=&for_exact=&for_none=&for_jt=&for_com=&for_loc=Denver%2C_CO&jtype=Full+Time&sort=relevance&limit=100&radius=0&postedDate=7&jtype=Full+Time&limit=100&radius=30&postedDate=7&jtype=Full+Time",
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
