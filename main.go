// crawler from mtg decks
package main

import (
	"log"
)

func main() {
	siteName := "mtgdecks"
	config, err := new(Config).loadConfig(siteName)
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	scraper := Scraper{
		config: config,
		name:   siteName,
	}

	scraper.url = "index-deck-list.html"
	scraper.loadDocument()
	scraper.fetchIndexPageData()

	scraper.url = "deck-page.html"
	scraper.loadDocument()
	scraper.fetchDeckPageData()

	scraper.url = "deck-with-side.html"
	scraper.loadDocument()
	scraper.fetchDeckPageData()
}
