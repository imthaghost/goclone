package crawler

import (
	"context"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strings"

	"github.com/gocolly/colly/v2"
)

// Collector searches for css, js, and images within a given link
// TODO improve for better performance
func Collector(ctx context.Context, url string, projectPath string, cookieJar *cookiejar.Jar, proxyString string, userAgent string, referer string, Depth int) error {
	log.Println(Depth)
	// create a new collector
	c := colly.NewCollector(colly.Async(true), colly.MaxDepth(Depth))
	setUpCollector(c, ctx, cookieJar, proxyString, userAgent)
	// search for all link tags that have a rel attribute that is equal to stylesheet - CSS
	c.OnHTML("link[rel='stylesheet']", func(e *colly.HTMLElement) {
		// hyperlink reference
		link := e.Attr("href")
		// extraction
		Extractor(url, e.Request.AbsoluteURL(link), projectPath, referer, userAgent, cookieJar)
	})

	// search for all script tags with src attribute -- JS
	c.OnHTML("script[src]", func(e *colly.HTMLElement) {
		// src attribute
		link := e.Attr("src")
		// extraction
		Extractor(url, e.Request.AbsoluteURL(link), projectPath, referer, userAgent, cookieJar)
	})

	// serach for all img tags with src attribute -- Images
	c.OnHTML("img[src]", func(e *colly.HTMLElement) {
		// src attribute
		link := e.Attr("src")
		if strings.HasPrefix(link, "data:image") || strings.HasPrefix(link, "blob:") {
			return
		}
		// Print link
		Extractor(url, e.Request.AbsoluteURL(link), projectPath, referer, userAgent, cookieJar)
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		Extractor(url, e.Request.AbsoluteURL(link), projectPath, referer, userAgent, cookieJar)
	})

	//Before making a request
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Referer", referer)
	})
	// Visit each url and wait for stuff to load :)
	if err := c.Visit(url); err != nil {
		return err
	}
	c.Wait()
	return nil
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

func setUpCollector(c *colly.Collector, ctx context.Context, cookieJar *cookiejar.Jar, proxyString, userAgent string) {

	if cookieJar != nil {
		c.SetCookieJar(cookieJar)
	}
	if proxyString != "" {
		c.SetProxy(proxyString)
	} else {
		c.WithTransport(cancelableTransport{ctx: ctx, transport: http.DefaultTransport})
	}
	if userAgent != "" {
		c.UserAgent = userAgent
	}
}
