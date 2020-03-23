package main

import (
	"fmt"
	"os"

	"github.com/imthaghost/goclone/crawler"
	"github.com/imthaghost/goclone/file"
	"github.com/imthaghost/goclone/html"
	"github.com/imthaghost/goclone/parser"
)

func main() {
	//fmt.Println(len(os.Args))
	// if the domain is valid and the url is valid return the text
	argsList := os.Args
	argslength := len(argsList)
	if argslength > 2 {
		fmt.Println("One argument for now")
		return
	} else if argslength <= 1 {
		fmt.Println("Help screen i think")
		return
	}
	url := os.Args[1]
	// grab the url from the
	if !parser.ValidateURL(url) && !parser.ValidateDomain(url) {
		fmt.Println("goclone <url>")
	} else if parser.ValidateDomain(url) {
		// use the domain as the project name
		name := url
		// CreateProject
		projectpath := file.CreateProject(name)
		// create the url
		validURL := parser.CreateURL(name)
		// Crawler
		crawler.Crawl(validURL, projectpath)
		// Restructure html
		html.LinkRestructure(projectpath)
	} else if parser.ValidateURL(url) {
		// get the hostname
		name := parser.GetDomain(url)
		// create project
		projectpath := file.CreateProject(name)
		// Crawler
		crawler.Crawl(url, projectpath)
		// Restructure html
		html.LinkRestructure(projectpath)
	} else {
		fmt.Print(url)
	}

}
