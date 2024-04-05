// crawler from mtg decks
package main

import (
	"log"
)

func main() {
	config, err := new(Config).loadConfig("mtgdecks")
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	scraper := Scraper{
		config: config,
		name:   "mtgdecks",
	}

	scraper.loadDocument("index-deck-list.html")
	scraper.fetchIndexPageData()

	scraper.loadDocument("deck-page.html")
	scraper.fetchDeckPageData()

	scraper.loadDocument("deck-with-side.html")
	scraper.fetchDeckPageData()
}
