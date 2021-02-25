package crawler

import (
	"fmt"

	"github.com/gocolly/colly"
)

// Collector searches for css, js, and images within a given link
// TODO improve for better performance
func Collector(url string, projectPath string) {
	// create a new collector
	c := colly.NewCollector(
		// asychronus boolean
		colly.Async(true),
	)

	// search for all link tags that have a rel attribute that is equal to stylesheet - CSS
	c.OnHTML("link[rel='stylesheet']", func(e *colly.HTMLElement) {
		// hyperlink reference
		link := e.Attr("href")
		// print css file was found
		fmt.Println("Css found", "-->", link)
		// extraction
		Extractor(e.Request.AbsoluteURL(link), projectPath)
	})

	// search for all script tags with src attribute -- JS
	c.OnHTML("script[src]", func(e *colly.HTMLElement) {
		// src attribute
		link := e.Attr("src")
		// Print link
		fmt.Println("Js found", "-->", link)
		// extraction
		Extractor(e.Request.AbsoluteURL(link), projectPath)
	})

	// serach for all img tags with src attribute -- Images
	c.OnHTML("img[src]", func(e *colly.HTMLElement) {
		// src attribute
		link := e.Attr("src")
		// Print link
		fmt.Println("Img found", "-->", link)
		// extraction
		Extractor(e.Request.AbsoluteURL(link), projectPath)
	})

	//Before making a request
	c.OnRequest(func(r *colly.Request) {
		link := r.URL.String()
		if url == link {
			HTMLExtractor(link, projectPath)
		}
	})

	// Visit each url and wait for stuff to load :)
	c.Visit(url)
	c.Wait()
}
