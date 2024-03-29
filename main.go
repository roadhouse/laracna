// crawler from mtg decks
package main

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/PuerkitoBio/goquery"
)

// Config the crawler
type Config struct {
	URL               string `toml:"url"`
	DeckPath          string `toml:"deck_path"`
	DeckDataPath      string `toml:"deck_data_path"`
	SideBoardPath     string `toml:"side_board_path"`
	SideBoardDataPath string `toml:"side_board_data_path"`
	OtherIndexPath    string `toml:"other_indexes_path"`
}

func loadConfig(filePath string) (map[string]Config, error) {
	var configs map[string]Config
	if _, err := toml.DecodeFile(filePath, &configs); err != nil {
		log.Printf("Error reading TOML file: %s", err)
		return nil, err
	}
	return configs, nil
}

func fetchPageContent(url string) (*goquery.Document, error) {
	doc, err := os.Open(url)
	// doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal("Error opening file:", err)
		return nil, err
	}
	// remove if use url
	defer doc.Close()

	content, err := goquery.NewDocumentFromReader(doc)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return content, nil
}

func main() {
	configs, err := loadConfig("config.toml")
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	// doc, err := fetchPageContent(configs["mtgdecks"].URL)
	doc, err := fetchPageContent("index-deck-list.html")
	if err != nil {
		log.Fatalf("Error parsing URL: %s", err)
	}

	scraper := Scraper{
		config: configs["mtgdecks"],
		name:   "mtgdecks",
		doc:    *doc,
		url:    configs["mtgdecks"].URL,
		// url: "index-page.html",
	}

	scraper.fetchIndexPageData()
	doc, err = fetchPageContent("deck-page.html")
	if err != nil {
		log.Fatalf("Error parsing URL: %s", err)
	}
	scraper.doc = *doc
	scraper.fetchDeckPageData()

	doc, err = fetchPageContent("deck-with-side.html")
	if err != nil {
		log.Fatalf("Error parsing URL: %s", err)
	}
	scraper.doc = *doc
	scraper.fetchDeckPageData()
}
