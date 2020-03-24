package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/imthaghost/goclone/auth"
	"github.com/imthaghost/goclone/crawler"
	"github.com/imthaghost/goclone/file"
	"github.com/imthaghost/goclone/flags"
	"github.com/imthaghost/goclone/html"
	"github.com/imthaghost/goclone/parser"
	"github.com/imthaghost/goclone/server"
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
	help, login, serve := flags.ParseFlags()
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
	} else if serve == true {
		url := os.Args[2]

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
			indexfile := projectpath + "/index.html"
			html.FormatHTML(indexfile)
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
			indexfile := projectpath + "/index.html"
			// Format html
			html.FormatHTML(indexfile)
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
