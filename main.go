package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/araddon/dateparse"
)

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

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

	nextDealUTC := doc.Find("span.js-clock").AttrOr("data-next-deal", "")
	nextDealParsed, err := dateparse.ParseStrict(nextDealUTC)
	nextDeal := nextDealParsed.In(time.Local)

	nextDealIn := time.Until(nextDeal)
	nextDealInFmt := fmtDuration(nextDealIn)

	fmt.Printf("\n    %s\n    %s\n\n", title, subtitle)
	fmt.Printf("Für CHF %s anstatt %s\n", price, originalPrice)
	fmt.Printf("Noch %s verfügbar\n", percentage)
	// Golang time formatting: https://flaviocopes.com/go-date-time-format/
	fmt.Printf("Nächster Deal am: %s (in %s)\n", nextDeal.Format("Mon Jan _2 15:04:05"), nextDealInFmt)
}
