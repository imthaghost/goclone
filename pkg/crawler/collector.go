package crawler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/cookiejar"

	"github.com/gocolly/colly"
)

// Collector searches for css, js, and images within a given link
// TODO improve for better performance
func Collector(ctx context.Context, url string, projectPath string, collyOpts ...func(*colly.Collector)) error {
	// create a new collector
	c := colly.NewCollector(append(collyOpts, colly.Async(true))...)
	c.WithTransport(cancelableTransport{ctx: ctx, transport: http.DefaultTransport})

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
	if err := c.Visit(url); err != nil {
		return err
	}
	c.Wait()
	return nil
}

// SetCookieJar returns a colly.Collector option that sets the cookie jar to the specified.
func SetCookieJar(jar *cookiejar.Jar) func(*colly.Collector) {
	return func(c *colly.Collector) { c.SetCookieJar(jar) }
}

type cancelableTransport struct {
	ctx       context.Context
	transport http.RoundTripper
}

func (t cancelableTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := t.ctx.Err(); err != nil {
		return nil, err
	}
	return t.transport.RoundTrip(req.WithContext(t.ctx))
}
