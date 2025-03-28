package html

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/PuerkitoBio/goquery"
)

// TODO: figure out what was done here at 4am
func arrange(projectDir string) error {
	indexfile := projectDir + "/index.html"
	input, err := os.ReadFile(indexfile)
	if err != nil {
		return err
	}

	// Create a single document from the entire HTML
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(input))
	if err != nil {
		return err
	}

	// Replace JS links in HTML
	doc.Find("script[src]").Each(func(i int, s *goquery.Selection) {
		data, exists := s.Attr("src")
		if exists {
			file := filepath.Base(data)
			newPath := "js/" + file
			s.SetAttr("src", newPath)
		}
	})

	// Replace CSS links in HTML
	doc.Find("link[rel='stylesheet']").Each(func(i int, s *goquery.Selection) {
		data, exists := s.Attr("href")
		if exists {
			file := filepath.Base(data)
			newPath := "css/" + file
			s.SetAttr("href", newPath)
		}
	})

	// Replace IMG links in HTML
	doc.Find("img[src]").Each(func(i int, s *goquery.Selection) {
		data, exists := s.Attr("src")
		if exists {
			file := filepath.Base(data)
			newPath := "imgs/" + file
			s.SetAttr("src", newPath)
		}
	})

	// Get the modified HTML
	html, err := doc.Html()
	if err != nil {
		return err
	}

	return os.WriteFile(indexfile, []byte(html), 0777)
}
