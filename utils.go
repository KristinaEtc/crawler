package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/asaskevich/govalidator"
)

const (
	httpStatusCodeBadURI      = "0"
	httpStatusCodeNetworkFail = "-1"

	notValidPrefix = "notValid"
)

func getURIsFromBodyRequest(bodyRaw []byte) []string {
	uris := make([]string, 0)
	uris = strings.Split(string(bodyRaw), "\n")

	log.Printf("[DEBUG] res=%+v, len=%d\n", uris, len(uris))
	for i, uri := range uris {
		valid := govalidator.IsRequestURL(uri)
		if !valid {
			uris[i] = notValidPrefix + uri
		}
	}

	log.Printf("[DEBUG] validated res=%+v, len=%d\n", uris, len(uris))
	return uris
}

func getMeta(resCh chan []string, uri string) {
	res := make([]string, 0)
	if strings.HasPrefix(uri, notValidPrefix) {
		resCh <- append(res, httpStatusCodeBadURI)
		return
	}

	doc, err := goquery.NewDocument(uri)
	if err != nil {
		log.Printf("[ERR] error processing uri=[%s]: %s\n", uri, err.Error())
		resCh <- append(res, httpStatusCodeNetworkFail)
		return
	}

	var metaDescription, metaKeywords, pageTitle, ogImage string

	pageTitle = doc.Find("title").Contents().Text()
	fmt.Printf("Page Title: '%s'\n", pageTitle)

	doc.Find("meta").Each(func(index int, item *goquery.Selection) {
		if item.AttrOr("name", "") == "description" {
			metaDescription = item.AttrOr("content", "")
		}
		if item.AttrOr("name", "") == "keywords" {
			metaKeywords = item.AttrOr("content", "")
		}
		metaKeywords = item.AttrOr("content", "")

		op, _ := item.Attr("property")
		con, _ := item.Attr("content")
		if op == "og:image" {
			ogImage = con
		}

	})
	fmt.Printf("Meta Description: '%s'\n", metaDescription)
	res = append(res, uri, pageTitle, metaDescription, metaKeywords, ogImage)

	resCh <- res
}
