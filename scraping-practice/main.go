package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

const scrapeBaseUrl = "https://www.nike.com/jp/launch?s=in-stock"

// const nextScrape = "https://www.nike.com/jp/launch/t/off-white-apparel-collection-fa21"

func ExampleScrape() {
	res, err := GetHttpHtmlContent(scrapeBaseUrl, "body")

	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(res))
	if err != nil {
		log.Fatal(err)
	}
	BasePageScrape(doc)

	doc.Find("aside ul li").Each(func(i int, s *goquery.Selection) {
		// fmt.Print(s.Html())
		//  ClapLv1TextBlock_Chie-TextBlock__Text--clamp2__1UeI0
		// title, exists := s.Find("a").Attr("href")
		title := s.Find("button[disabled!='']").Text()
		fmt.Println(title)
		// fmt.Printf("Review %d: %s\n", i, title)
	})
}

const endpoint = "https://discord.com/api/webhooks/860020793588449312/GKiCb6aDqQKoYnLy7Smishignvd4cT6lnVmOYoNWGkheRAFh00k_MpfusWRWxlAtYcoD"

type Discord struct {
	Content string `json:"content"`
}

func SendDiscord(link string) {
	reqBody := &Discord{
		Content: link,
	}
	jsonString, err := json.Marshal(reqBody)
	if err != nil {
		panic("Cannot Send for reason json")
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonString))
	if err != nil {
		panic("Error: request")
	}

	req.Header.Set("Content-Type", "application/json")
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		panic("Error")
	}
	defer resp.Body.Close()

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic("Error")
	}

	fmt.Printf("%#v", string(byteArray))
}

func BasePageScrape(doc *goquery.Document) {
	doc.Find("figure").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Find("a").Attr("href")
		SendDiscord(link)
		fmt.Println(link)
	})
}

func main() {
	ExampleScrape()
}
