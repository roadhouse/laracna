package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"github.com/PuerkitoBio/goquery"
)

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
  href, _ := doc.Attr("href")

  return href
}

func fetchDeckList(i int, s *goquery.Selection) string {
  return strings.TrimSpace(s.Text())
}

func fetchDeckUrls(doc *goquery.Document) {
  xpath := "div.decks div div table tbody tr td a"
  deckList := doc.Find(xpath).Map(fetchLink)

  fmt.Println(deckList)
}

func fetchDeckData(doc *goquery.Document) {
  xpath := "div.wholeDeck td.number"
  j := doc.Find(xpath).Map(fetchDeckList)
  myMap := make(map[string]int)
  for _, line := range j {
    o := strings.Split(line, "\n")
    quantity, _ := strconv.Atoi(o[0])
    myMap[strings.TrimSpace(o[1])] = quantity
  }

  fmt.Println(myMap)
}

func fetchSideboarData(doc *goquery.Document) {
  myMap := make(map[string]int)
  sideBoardPath := "div.wholeDeck div.row div:last-child .Sideboard"
  hasSideboard := doc.Find(sideBoardPath).Length() != 0
  if hasSideboard {
    xpath := "div.wholeDeck div.row div:last-child td.number"
    j := doc.Find(xpath).Map(fetchDeckList)
    for _, line := range j {
      o := strings.Split(line, "\n")
      quantity, _ := strconv.Atoi(o[0])
      myMap[strings.TrimSpace(o[1])] = quantity
    }
  }

  fmt.Println(myMap)
}

func fetchOtherIndexLinks(doc *goquery.Document) {
  xpath := "ul.pagination li a"
  indexesList := doc.Find(xpath).Map(fetchLink)

  fmt.Println(indexesList)
}

func fetchIndexPageData(url string) {
  indexPage, err := fetchPageContent(url)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println("Links")
  fetchOtherIndexLinks(indexPage)
  fetchDeckUrls(indexPage)
}

func fetchDeckPageData(url string) {
  deckPage, err := fetchPageContent(url)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println("Deck")
  fetchDeckData(deckPage)
  fmt.Println("Sideboard")
  fetchSideboarData(deckPage)
}

func main() {
  fetchIndexPageData("index-deck-list.html")
  fetchDeckPageData("deck-page.html")
  fetchDeckPageData("deck-with-side.html")
}
