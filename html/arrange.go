package html

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// arrangeJS arranges javascript in index file
func arrangeJS(projectDir string) {
	// css project directory
	indexfile := projectDir + "/index.html"
	input, err := ioutil.ReadFile(indexfile)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(input), "\n")

	for index, line := range lines {
		b := []byte(line)
		r := bytes.NewReader(b)
		doc, err := goquery.NewDocumentFromReader(r)
		if err != nil {
			panic(err)
		}
		doc.Find("script[src]").Each(func(i int, s *goquery.Selection) {
			data, exists := s.Attr("src")
			if exists {
				file := filepath.Base(data)
				s.SetAttr("src", "js/"+file)
				data, exists := s.Attr("src")
				lines[index] = fmt.Sprintf(`<script src="%s"></script>`, data)
				if exists {
				}
			}
		})
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(indexfile, []byte(output), 0777)
	if err != nil {
		log.Fatalln(err)
	}
}

// arrangeCSS arrages css files in index file
func arrangeCSS(projectDir string) {
	indexfile := projectDir + "/index.html"
	input, err := ioutil.ReadFile(indexfile)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(input), "\n")
	// uh oh :(
	if err != nil {
		panic(err)
	}
	for index, line := range lines {
		b := []byte(line)
		r := bytes.NewReader(b)
		doc, err := goquery.NewDocumentFromReader(r)
		if err != nil {
			panic(err)
		}
		// Find where link has a rel attribute equal to stylesheets
		doc.Find("link[rel='stylesheet']").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the hyperlink reference
			data, exists := s.Attr("href")
			if exists {
				file := filepath.Base(data)
				s.SetAttr("href", "css/"+file)
				data, exists := s.Attr("href")
				lines[index] = fmt.Sprintf(`<link rel="stylesheet" type="text/css" href="%s">`, data)
				if exists {
				}
			}
		})
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(indexfile, []byte(output), 0777)
	if err != nil {
		log.Fatalln(err)
	}
}

// arrangeImgs arranges all images in project
func arrangeImgs(projectDir string) {
	// css project directory
	indexfile := projectDir + "/index.html"
	input, err := ioutil.ReadFile(indexfile)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(input), "\n")
	// images, err := ioutil.ReadDir(projectDir + "/imgs")
	// uh oh :(
	if err != nil {
		panic(err)
	}
	for index, line := range lines {
		b := []byte(line)
		r := bytes.NewReader(b)
		doc, err := goquery.NewDocumentFromReader(r)
		if err != nil {
			panic(err)
		}

		doc.Find("img[src]").Each(func(i int, s *goquery.Selection) {

			data, exists := s.Attr("src")
			if exists {
				file := filepath.Base(data)
				s.SetAttr("src", "imgs/"+file)
				data, exists := s.Attr("src")
				lines[index] = fmt.Sprintf(`<img src="%s">`, data)
				if exists {
				}
			}
		})
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(indexfile, []byte(output), 0777)
	if err != nil {
		log.Fatalln(err)
	}
}
