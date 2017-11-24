package main

import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
)

func byteToSlice(rawURIs []byte) []string {
	return nil
}

func getMeta(resCh chan []byte, uri string) {
	doc, err := goquery.NewDocument(uri)
	if err != nil {
		log.Printf("[ERR] error processing uri=[%s]: %s\n", uri, err.Error())
		return
	}

	var metaDescription string
	var pageTitle string

	// use CSS selector found with the browser inspector
	// for each, use index and item
	pageTitle = doc.Find("title").Contents().Text()

	doc.Find("meta").Each(func(index int, item *goquery.Selection) {
		if item.AttrOr("name", "") == "description" {
			metaDescription = item.AttrOr("content", "")
		}
	})
	fmt.Printf("Page Title: '%s'\n", pageTitle)
	fmt.Printf("Meta Description: '%s'\n", metaDescription)

	res := []byte("test")
	resCh <- res
}
