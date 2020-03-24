package cmd

import (
	"fmt"
	"os/exec"

	"github.com/imthaghost/goclone/auth"
	"github.com/imthaghost/goclone/crawler"
	"github.com/imthaghost/goclone/file"
	"github.com/imthaghost/goclone/html"
	"github.com/imthaghost/goclone/parser"
	"github.com/imthaghost/goclone/server"
)

// Clone the given site :)
func cloneSite(args []string) {
	url := args[0]
	if Login == true {
		// grab the url from the
		if !parser.ValidateURL(url) && !parser.ValidateDomain(url) {
			fmt.Println("goclone <url>")
		} else if parser.ValidateDomain(url) {
			// grab user credentials
			username, password := auth.Credentials()
			// use the domain as the project name
			name := url
			// CreateProject
			projectpath := file.CreateProject(name)
			// create the url
			validURL := parser.CreateURL(name)
			// Crawler
			crawler.LoginCollector(projectpath, validURL, username, password)
			// Restructure html
			html.LinkRestructure(projectpath)

		} else if parser.ValidateURL(url) {
			// grab user credentials
			username, password := auth.Credentials()
			// get the hostname
			name := parser.GetDomain(url)
			// CreateProject
			projectpath := file.CreateProject(name)
			// Crawler
			crawler.LoginCollector(projectpath, url, username, password)
			// Restructure html
			html.LinkRestructure(projectpath)

		} else {
			fmt.Print(url)
		}
	} else if Serve == true {
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
			err := exec.Command("open", "http://localhost:5000").Start()
			if err != nil {
				panic(err)
			}
			server.Serve(projectpath)

		} else if parser.ValidateURL(url) {
			// get the hostname
			name := parser.GetDomain(url)
			// create project
			projectpath := file.CreateProject(name)
			// Crawler
			crawler.Crawl(url, projectpath)
			// Restructure html
			html.LinkRestructure(projectpath)
			err := exec.Command("open", "http://localhost:5000").Start()
			if err != nil {
				panic(err)
			}
			server.Serve(projectpath)
		} else {
			fmt.Print(url)
		}
	} else {
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
			err := exec.Command("open", projectpath+"/index.html").Start()
			if err != nil {
				panic(err)
			}
		} else if parser.ValidateURL(url) {
			// get the hostname
			name := parser.GetDomain(url)
			// create project
			projectpath := file.CreateProject(name)
			// Crawler
			crawler.Crawl(url, projectpath)
			// Restructure html
			html.LinkRestructure(projectpath)
			err := exec.Command("open", projectpath+"/index.html").Start()
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Print(url)
		}
	}
}
