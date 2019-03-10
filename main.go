package main

import (
	"flag"
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

func dealPrice(doc *goquery.Document) string {
	return doc.Find(".product-pricing__prices-new-price").First().Text()
}

func origPrice(doc *goquery.Document) string {
	originalPrice := doc.Find("strong.product-pricing__prices-old-price").First().Text()
	originalPrice = strings.TrimSpace(originalPrice)
	originalPrice = strings.TrimSuffix(originalPrice, "*")
	originalPrice = strings.TrimSpace(originalPrice)
	return originalPrice
}

func dealPercentage(doc *goquery.Document) string {
	return doc.Find(".product-progress__availability").First().Text()
}

func main() {
	outputAllInfo := false
	dealAvailabilityFlg := flag.Bool("availability", false, "Availability")
	dealPriceFlg := flag.Bool("price", false, "Price")
	dealNameFlg := flag.Bool("name", false, "Name")
	dealTitleFlg := flag.Bool("title", false, "Title")
	dealSubtitleFlg := flag.Bool("subtitle", false, "Subtitle")
	dealURLFlg := flag.String("url", "default", "Deal url. So far supported: 'https://daydeal.ch', 'https://www.daydeal.ch/deal-of-the-week', 'https://blickdeal.ch'")

	flag.Parse()
	if (flag.NFlag() == 1 && *dealURLFlg != "default") || flag.NFlag() == 0 {
		outputAllInfo = true
	}

	//If the user inputs no url we need to ensure that it's pointing to daydeal.ch
	if *dealURLFlg == "default" {
		*dealURLFlg = "https://daydeal.ch"
	}

	doc, err := goquery.NewDocument(*dealURLFlg)
	if err != nil {
		log.Fatal(err)
	}

	title := dealName(doc)
	subtitle := dealPrecName(doc)

	price := dealPrice(doc)
	originalPrice := origPrice(doc)

	percentage := dealPercentage(doc)

	if *dealTitleFlg == true {
		fmt.Printf("%s\n", title)
	}
	if *dealNameFlg == true || outputAllInfo == true {
		fmt.Printf("\n    %s\n    %s\n\n", title, subtitle)
	}
	if *dealSubtitleFlg == true {
		fmt.Printf("%s\n", subtitle)
	}

	if *dealPriceFlg == true || outputAllInfo == true {
		fmt.Printf("Für CHF %s anstatt %s\n", price, originalPrice)
	}
	if *dealAvailabilityFlg == true || outputAllInfo == true {
		fmt.Printf("Noch %s verfügbar\n", percentage)
	}
}
