package crawler

import (
	"fmt"
	"os"

	"github.com/gocolly/colly"
)

// LoginCollector ...
func LoginCollector(projectPath string, site string, username string, password string) {

	// csrfTokenSelector := "#main-container > section.content > main > div > div.auth-form.sign-in-form > form > input[type=hidden]:nth-child(2)"
	//CSRF Token Authorization
	// c.OnHTML(csrfTokenSelector, func(e *colly.HTMLElement) {
	// 	//authenticity token
	//token := e.Attr("value")
	// })
	// new collector
	c := colly.NewCollector()
	//authenticate
	err := c.Post("https://www.makeschool.com/login", map[string]string{"user[email]": username, "user[password]": password})
	if err != nil {
		panic(err)
	}

	// fmt.Print(c.Cookies("https://www.makeschool.com/dashboard"))
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
		r.Ctx.Put("url", link)

	})
	//Response of each visited page
	c.OnResponse(func(r *colly.Response) {
		// link := r.Ctx.Get("url")
		// check if the url being visited is the root for searching if so write it as a page

		body := r.Body
		f, err := os.OpenFile(projectPath+"/"+"index.html", os.O_RDWR|os.O_CREATE, 0777)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		f.Write(body)

	})

	// start scraping
	c.Visit("https://www.makeschool.com/dashboard")

}
