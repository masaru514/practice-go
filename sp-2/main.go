package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
)

// func inPageAccess(url string) {
// 	c := colly.NewCollector()
// 	c.OnHTML()
// 	c.Visit(url)
// }

func main() {
	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// crete a timeout
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	start := time.Now()

	var res string
	err := chromedp.Run(ctx,
		emulation.SetUserAgentOverride("WebScraper 1.0"),
		chromedp.Navigate(`https://chiebukuro.yahoo.co.jp`),
		// wait for footer element is visible (ie, page is loaded)
		chromedp.ScrollIntoView(`footer`),
		// chromedp.WaistVisible(`footer < div`),
		chromedp.Text(`h2`, &res, chromedp.NodeVisible, chromedp.ByQuery),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("h1 contains: '%s'\n", res)
	fmt.Printf("\n\nTook: %f secs\n", time.Since(start).Seconds())

	// c := colly.NewCollector()

	// colly.AllowedDomains("www.nike.com")

	// c.OnHTML(".size-layout", func(e *colly.HTMLElement) {
	// 	link := e.Attr("href")
	// 	// Print link
	// 	fmt.Printf("Link found: %q -> %s\n", e.DOM, link)
	// 	// Visit link found on page
	// 	// Only those links are visited which are in AllowedDomains
	// 	// c.Visit(e.Request.AbsoluteURL(link))
	// })

	// c.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("%v", r.URL.String())
	// })

	// // アクセス
	// c.Visit("https://www.nike.com/jp/launch/t/nike-sb-parra-japan-federation-kit")
}
