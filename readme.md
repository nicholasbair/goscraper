Scraper Goals
[x] Scrape 1 page of Indeed
[x] Paginate through results
[x] Add support for Dice (scraper functions should be generic)
[ ] Add channels/goroutines
[ ] Expose data as API
[ ] Add database?
[ ] Call node microservice (LinkedIn, Muse) or re-write in golang

API things
- localhost/indeed/jobs?keywords=stuff&location=places
  - Need to parse keywords, location, provider etc from request
  - Need function to build initial & page URL for each provider
- How to handle pagination?
- Encode response as JSON

Server receives a request, parse the URL, create a REQ and a channel, call scrape with the REQ.  Scrape will then look up the right config, and call the relevant scraper functions.  Once slice of results urls is built, each call to getJobData should be with a new channel and goroutine

API things to handle later (much later)
- Authenication# goscraper
