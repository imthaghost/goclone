package goclone

import (
	"fmt"

	"os/exec"

	"github.com/imthaghost/goclone/pkg/crawler"
	"github.com/imthaghost/goclone/pkg/file"
	"github.com/imthaghost/goclone/pkg/html"
	"github.com/imthaghost/goclone/pkg/parser"
	"github.com/imthaghost/goclone/pkg/server"
)

// Clone the given site :)
func cloneSite(args []string) {
	url := args[0]

	if Serve == true {
		// grab the url from the
		if !parser.ValidateURL(url) && !parser.ValidateDomain(url) {
			fmt.Println("goclone <url>")
		} else if parser.ValidateDomain(url) {
			// use the domain as the project name
			name := url
			// CreateProject
			projectPath := file.CreateProject(name)
			// create the url
			validURL := parser.CreateURL(name)
			// Crawler
			crawler.Crawl(validURL, projectPath)
			// Restructure html
			html.LinkRestructure(projectPath)
			err := exec.Command("open", "http://localhost:5000").Start()
			if err != nil {
				panic(err)
			}
			server.Serve(projectPath)

		} else if parser.ValidateURL(url) {
			// get the hostname
			name := parser.GetDomain(url)
			// create project
			projectPath := file.CreateProject(name)
			// Crawler
			crawler.Crawl(url, projectPath)
			// Restructure html
			html.LinkRestructure(projectPath)
			err := exec.Command("open", "http://localhost:5000").Start()
			if err != nil {
				panic(err)
			}
			server.Serve(projectPath)
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
			projectPath := file.CreateProject(name)
			// create the url
			validURL := parser.CreateURL(name)
			// Crawler
			crawler.Crawl(validURL, projectPath)
			// Restructure html
			html.LinkRestructure(projectPath)
			if Open {
				// automatically open project
				err := exec.Command("open", projectPath+"/index.html").Start()
				if err != nil {
					panic(err)
				}
			}

		} else if parser.ValidateURL(url) {
			// get the hostname
			name := parser.GetDomain(url)
			// create project
			projectPath := file.CreateProject(name)
			// Crawler
			crawler.Crawl(url, projectPath)
			// Restructure html
			html.LinkRestructure(projectPath)
			if Open {
				err := exec.Command("open", projectPath+"/index.html").Start()
				if err != nil {
					panic(err)
				}
			}
		} else {
			fmt.Print(url)
		}
	}
}
