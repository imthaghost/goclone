package main

import (
	"fmt"
	"os"

	"github.com/imthaghost/goclone/auth"
	"github.com/imthaghost/goclone/crawler"
	"github.com/imthaghost/goclone/file"
	"github.com/imthaghost/goclone/flags"
	"github.com/imthaghost/goclone/html"
	"github.com/imthaghost/goclone/parser"
)

func main() {
	usage := `

		Usage:
		goclone
		goclone <url>
		goclone --help
		goclone --version
		goclone --v
	
	Options:
		<url>  Optional url argument.
		--help  Show help screen.
		--version  Show version.
		--v		Verbose output`
	argsList := os.Args
	argslength := len(argsList)
	if argslength <= 1 {
		fmt.Println(usage)
		return
	}
	help, login := flags.ParseFlags()
	if help {
		fmt.Println(usage)
		return
	}
	if login == true {
		url := os.Args[2]
		// grab the url from the
		if !parser.ValidateURL(url) && !parser.ValidateDomain(url) {
			fmt.Println("goclone <url>")
		} else if parser.ValidateDomain(url) {
			username, password := auth.Credentials()
			fmt.Printf("Username: %s, Password: %s\n", username, password)
			// use the domain as the project name
			name := url
			// CreateProject
			file.CreateProject(name)
			// create the url
			validURL := parser.CreateURL(name)
			// Crawler
			crawler.LoginCollector(validURL)

		} else if parser.ValidateURL(url) {
			username, password := auth.Credentials()
			fmt.Printf("Username: %s, Password: %s\n", username, password)
			// get the hostname
			name := parser.GetDomain(url)
			// create project
			file.CreateProject(name)
			// Crawler
			crawler.LoginCollector(url)

		} else {
			fmt.Print(url)
		}
	} else {
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

}
