package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Scraper struct {
	config Config
	name   string
	doc    goquery.Document
	url    string
}

func (scraper *Scraper) loadDocument(url string) (Scraper, error) {
	doc, err := os.Open(url)
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

func (scraper *Scraper) fetchDeckUrls() {
	deckList := scraper.doc.Find(scraper.config.DeckPath).Map(scraper.fetchLink)

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
	j := scraper.doc.Find(scraper.config.DeckDataPath).Map(scraper.fetchDeckList)
	myMap := scraper.quantityAndCardName(j)
	fmt.Println(myMap)
}

func (scraper *Scraper) fetchSideboarData() {
	myMap := make(map[string]int)
	hasSideboard := scraper.doc.Find(scraper.config.SideBoardPath).Length() != 0
	if hasSideboard {
		j := scraper.doc.Find(scraper.config.SideBoardDataPath).Map(scraper.fetchDeckList)
		myMap = scraper.quantityAndCardName(j)
	}
	fmt.Println(myMap)
}

func (scraper *Scraper) fetchOtherIndexLinks() {
	indexesList := scraper.doc.Find(scraper.config.OtherIndexPath).Map(scraper.fetchLink)

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

func (scraper *Scraper) fetchLink(i int, doc *goquery.Selection) string {
	_ = i
	href, _ := doc.Attr("href")

	return href
}

func (scraper *Scraper) fetchDeckList(i int, s *goquery.Selection) string {
	_ = i
	return strings.TrimSpace(s.Text())
}
