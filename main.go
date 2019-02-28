package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	doc, err := goquery.NewDocument("https://www.daydeal.ch/")
	if err != nil {
		log.Fatal(err)
	}

	title := doc.Find(".product-description__title1").First().Text()
	subtitle := doc.Find(".product-description__title2").First().Text()

	price := doc.Find(".product-pricing__prices-new-price").First().Text()
	originalPrice := doc.Find("div.originalPrice span").First().Text()
	originalPrice = strings.TrimSpace(originalPrice)
	originalPrice = strings.TrimSuffix(originalPrice, "*")
	originalPrice = strings.TrimSpace(originalPrice)

	percentage := doc.Find(".product-progress__availability").First().Text()

	nextDeal := doc.Find("span.js-clock").AttrOr("data-next-deal", "")

	fmt.Printf("\n    %s\n    %s\n\n", title, subtitle)
	fmt.Printf("Für CHF %s anstatt %s\n", price, originalPrice)
	fmt.Printf("Noch %s verfügbar\n", percentage)
	fmt.Printf("Nächster Deal am: %s\n", nextDeal)
}
