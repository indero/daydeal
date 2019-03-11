package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/goware/urlx"
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

func dealNext(doc *goquery.Document) time.Time {
	// the website returns the moment of the next deal like "2006-01-02 15:04:05"
	// in UTC without announcingt that it is UTC.
	nextDealUTC := doc.Find("span.js-clock").AttrOr("data-next-deal", "")
	nextDeal, _ := time.ParseInLocation("2006-01-02 15:04:05", nextDealUTC, time.Local)
	return nextDeal
}

func sanitizeURL(inputURL string) string {
	url, _ := urlx.Parse(inputURL)
	url.Scheme = "https"
	normalizedURL, _ := urlx.Normalize(url)
	return normalizedURL
}

func main() {
	outputAllInfo := false
	defaultURL := "https://www.daydeal.ch"

	dealAvailabilityFlg := flag.Bool("availability", false, "Availability")
	dealPriceFlg := flag.Bool("price", false, "Price")
	dealNameFlg := flag.Bool("name", false, "Name")
	dealTitleFlg := flag.Bool("title", false, "Title")
	dealSubtitleFlg := flag.Bool("subtitle", false, "Subtitle")
	dealURLFlg := flag.String("url", "default", "Deal url. So far supported: 'https://daydeal.ch', 'https://www.daydeal.ch/deal-of-the-week', 'https://blickdeal.ch'")
	dealNextFlg := flag.Bool("next", false, "Show time and date of the next deal")

	flag.Parse()
	if (flag.NFlag() == 1 && *dealURLFlg != "default") || flag.NFlag() == 0 {
		outputAllInfo = true
	}

	//If the user inputs no url we need to ensure that it's pointing to daydeal.ch
	if *dealURLFlg == "default" {
		*dealURLFlg = defaultURL
	}

	normalizedURL := sanitizeURL(*dealURLFlg)

	doc, err := goquery.NewDocument(normalizedURL)
	if err != nil {
		log.Fatal(err)
	}

	title := dealName(doc)
	subtitle := dealPrecName(doc)

	price := dealPrice(doc)
	originalPrice := origPrice(doc)

	percentage := dealPercentage(doc)

	nextDeal := dealNext(doc)

	nextDealIn := time.Until(nextDeal)
	nextDealInFmt := fmtDuration(nextDealIn)

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

	if *dealNextFlg == true || outputAllInfo == true {
		// Golang time formatting: https://flaviocopes.com/go-date-time-format/
		fmt.Printf("Nächster Deal am: %s (in %s)\n", nextDeal.Format("Mon Jan _2 15:04:05"), nextDealInFmt)
	}

	if outputAllInfo == true {
		fmt.Printf("%s\n", normalizedURL)
	}
}
