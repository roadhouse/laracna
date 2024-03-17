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

func fetchDeckUrl(i int, doc *goquery.Selection) string {
  href, _ := doc.Attr("href")

  return href
}

func fetchDeckList(i int, s *goquery.Selection) string {
    return strings.TrimSpace(s.Text())
}

func main() {
  // ------------------- List decks urls
  doc, err := fetchPageContent("index-deck-list.html")
  if err != nil {
    log.Fatal(err)
  }
  xpath := "div.decks div div table tbody a"
  deckList := doc.Find(xpath).Map(fetchDeckUrl)
  fmt.Println(deckList)

  // ------------------- Parsing deck 
  doc, err = fetchPageContent("deck-page.html")
  if err != nil {
    log.Fatal(err)
  }
  xpath = "div.wholeDeck td.number"
  j := doc.Find(xpath).Map(fetchDeckList)
  myMap := make(map[string]int)
  for _, line := range j {
    o := strings.Split(line, "\n")
    quantity, _ := strconv.Atoi(o[0])
    myMap[strings.TrimSpace(o[1])] = quantity
  }
  fmt.Println(myMap)
}
