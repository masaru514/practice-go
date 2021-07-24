package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
)

func GetHttpHtmlContent(url string, sel interface{}) (string, error) {
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true), // debug usage
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
	}

	//Initialization parameters, first pass an empty data
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)

	c, _ := chromedp.NewExecAllocator(context.Background(), options...)

	// create context
	chromeCtx, _ := chromedp.NewContext(c, chromedp.WithLogf(log.Printf))

	// Execute an empty task to create a Chrome instance in advance
	chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...)

	//Create a context with a timeout of 40s
	timeoutCtx, cancel := context.WithTimeout(chromeCtx, 40*time.Second)
	defer cancel()

	var htmlContent string
	err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(ctx context.Context) error {
			node, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}
			res, er := dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
			htmlContent = res
			return er
		}),
	)
	if err != nil {
		log.Fatalf("Run err : %v\n", err)
		return "", err
	}
	return htmlContent, nil
}

func ExampleScrape() {
	res, err := GetHttpHtmlContent("https://chiebukuro.yahoo.co.jp/", "body")

	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v", doc.Find("#all_rnk"))

	doc.Find("#all_rnk").Each(func(i int, s *goquery.Selection) {
		fmt.Printf("%v", s)
		//  ClapLv1TextBlock_Chie-TextBlock__Text--clamp2__1UeI0
		title := s.Find("div").Text()
		fmt.Printf("Review %d: %s\n", i, title)
	})
}

func main() {
	ExampleScrape()
}
