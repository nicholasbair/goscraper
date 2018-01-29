package goscraper

import (
	"strings"
)

func checkError(err error) {
	if err != nil {
		ce <- err
		close(ce)
	}
}

func cleanString(s string) []string {
	return strings.Split(strings.TrimSpace(strings.Replace(s, ",", "", -1)), " ")
}

func strip(s string) string {
	x := strings.Replace(s, "\t", "", -1)
	x = strings.Replace(x, "\n", "", -1)
	return x
}

func cleanMap(c Config, p map[string][]string) {
	for k := range c.QueryMap {
		if _, ok := p[k]; !ok {
			delete(p, k)
		}
	}
}
