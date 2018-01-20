package goscraper

import "strings"

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func cleanString(s string) []string {
	newS := strings.Replace(s, ",", "", -1)
	return strings.Split(strings.TrimSpace(newS), " ")
}
