package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func dealName(doc *goquery.Document) string {
	return doc.Find(".product-description__title1").First().Text()
}

func dealPrecName(doc *goquery.Document) string {
	return doc.Find(".product-description__title2").First().Text()
}

func origPrice(doc *goquery.Document) string {
	originalPrice := doc.Find("strong.product-pricing__prices-old-price").First().Text()
	originalPrice = strings.TrimSpace(originalPrice)
	originalPrice = strings.TrimSuffix(originalPrice, "*")
	originalPrice = strings.TrimSpace(originalPrice)
	return originalPrice
}

func main() {
	doc, err := goquery.NewDocument("https://www.daydeal.ch/")
	if err != nil {
		log.Fatal(err)
	}

	title := dealName(doc)
	subtitle := dealPrecName(doc)

	price := doc.Find(".product-pricing__prices-new-price").First().Text()
	originalPrice := origPrice(doc)

	percentage := doc.Find(".product-progress__availability").First().Text()

	fmt.Printf("\n    %s\n    %s\n\n", title, subtitle)
	fmt.Printf("Für CHF %s anstatt %s\n", price, originalPrice)
	fmt.Printf("Noch %s verfügbar\n", percentage)
}
