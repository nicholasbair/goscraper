package goscraper

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
