// crawler from mtgdecks.com
package main

import (
	"fmt"
	"log"
)

func main() {
	scraper, err := new(Scraper).LoadScraper("mtgdecks")
	if err != nil {
		log.Fatalf("Error loading scraper: %s", err)
	}

	index := scraper.ExtractIndexPageData("index-deck-list.html")
	fmt.Printf("index page data: %+v\n\n", index)

	deck := scraper.ExtractDeckInfo("deck-with-side.html")
	fmt.Printf("deck with sideboard: %+v\n\n", deck)

	deck = scraper.ExtractDeckInfo("deck-page.html")
	fmt.Printf("deck without sideboard: %+v\n\n", deck)
}
