package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type IndexPage struct {
	OtherIndexes []string
	DeckUrls     []string
}

type Deck struct {
	DeckName  string
	Main      map[string]int
	Sideboard map[string]int
}

type Scraper struct {
	config Config
	name   string
	doc    goquery.Document
	url    string
}

func (scraper *Scraper) loadDocument() (Scraper, error) {
	doc, err := os.Open(scraper.url)
	// doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal("Error opening file:", err)
		return Scraper{}, err
	}
	// remove if use url
	defer doc.Close()

	content, err := goquery.NewDocumentFromReader(doc)
	if err != nil {
		log.Fatal(err)
		return Scraper{}, err
	}

	scraper.doc = *content
	return *scraper, nil
}

func (scraper *Scraper) fetchDeckUrls() []string {
	return scraper.doc.Find(scraper.config.DeckPath).Map(scraper.fetchLink)
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

func (scraper *Scraper) fetchDeckData() map[string]int {
	j := scraper.doc.Find(scraper.config.DeckDataPath).Map(scraper.fetchDeckList)
	myMap := scraper.quantityAndCardName(j)
	return myMap
}

func (scraper *Scraper) fetchSideboarData() map[string]int {
	myMap := make(map[string]int)
	hasSideboard := scraper.doc.Find(scraper.config.SideBoardPath).Length() != 0
	if hasSideboard {
		j := scraper.doc.Find(scraper.config.SideBoardDataPath).Map(scraper.fetchDeckList)
		myMap = scraper.quantityAndCardName(j)
	}
	return myMap
}

func (scraper *Scraper) fetchOtherIndexLinks() []string {
	return scraper.doc.Find(scraper.config.OtherIndexPath).Map(scraper.fetchLink)
}

func (scraper *Scraper) ExtractIndexPageData(url string) IndexPage {
	scraper.url = url
	scraper.loadDocument()

	return IndexPage{
		OtherIndexes: scraper.fetchOtherIndexLinks(),
		DeckUrls:     scraper.fetchDeckUrls(),
	}
}

func (scraper *Scraper) fetchLink(i int, doc *goquery.Selection) string {
	_ = i
	href, _ := doc.Attr("href")

	return href
}

func (scraper *Scraper) fetchDeckList(i int, s *goquery.Selection) string {
	_ = i
	return strings.TrimSpace(s.Text())
}

func (scraper *Scraper) fetchDeckName() string {
	return scraper.doc.Find(scraper.config.DeckNamePath).Text()
}

func (scraper *Scraper) ExtractDeckInfo(url string) Deck {
	scraper.url = url
	scraper.loadDocument()

	return Deck{
		DeckName:  scraper.fetchDeckName(),
		Main:      scraper.fetchDeckData(),
		Sideboard: scraper.fetchSideboarData(),
	}
}

func (scraper *Scraper) LoadScraper(configName string) (Scraper, error) {
	config, err := new(Config).loadConfig(configName)
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
		return Scraper{}, err
	}

	return Scraper{
		config: config,
		name:   configName,
	}, nil
}
