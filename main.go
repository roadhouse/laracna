// crawler from mtg decks
package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

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

type Scraper struct {
	config Config
	name   string
	doc    goquery.Document
	url    string
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

func fetchLink(i int, doc *goquery.Selection) string {
	_ = i
	href, _ := doc.Attr("href")

	return href
}

func fetchDeckList(i int, s *goquery.Selection) string {
	_ = i
	return strings.TrimSpace(s.Text())
}

func (scraper *Scraper) fetchDeckUrls() {
	deckList := scraper.doc.Find(scraper.config.DeckPath).Map(fetchLink)

	fmt.Println(deckList)
}

func (scraper *Scraper) quantityAndCardName(entries []string) map[string]int {
	myMap := make(map[string]int)
	for _, line := range entries {
		o := strings.Split(line, "\n")
		quantity, _ := strconv.Atoi(o[0])
		myMap[strings.TrimSpace(o[1])] = quantity
	}
	return myMap
}

func (scraper *Scraper) fetchDeckData() {
	j := scraper.doc.Find(scraper.config.DeckDataPath).Map(fetchDeckList)
	myMap := scraper.quantityAndCardName(j)
	fmt.Println(myMap)
}

func (scraper *Scraper) fetchSideboarData() {
	myMap := make(map[string]int)
	hasSideboard := scraper.doc.Find(scraper.config.SideBoardPath).Length() != 0
	if hasSideboard {
		j := scraper.doc.Find(scraper.config.SideBoardDataPath).Map(fetchDeckList)
		myMap = scraper.quantityAndCardName(j)
	}
	fmt.Println(myMap)
}

func (scraper *Scraper) fetchOtherIndexLinks() {
	indexesList := scraper.doc.Find(scraper.config.OtherIndexPath).Map(fetchLink)

	fmt.Println(indexesList)
}

func (scraper *Scraper) fetchIndexPageData() {
	scraper.fetchOtherIndexLinks()
	scraper.fetchDeckUrls()
}

func (scraper *Scraper) fetchDeckPageData() {
	fmt.Println("Deck")
	scraper.fetchDeckData()
	fmt.Println("Sideboard")
	scraper.fetchSideboarData()
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
