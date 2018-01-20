package goscraper

import "strings"

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func cleanString(s string) []string {
	return strings.Split(strings.TrimSpace(strings.Replace(s, ",", "", -1)), " ")
}

func cleanMap(c Config, p map[string][]string) {
	for k := range c.QueryMap {
		if _, ok := p[k]; !ok {
			delete(p, k)
		}
	}
}
