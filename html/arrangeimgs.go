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

// ArrangeImgs arranges all images in project
func ArrangeImgs(projectDir string) {
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
