package html

import (
	"regexp"
)

// CSSFinder takes in a byte slice and search through HTML data for css files
func CSSFinder(htmlData []byte) []string {
	re := regexp.MustCompile(`href="([^"]+\.css)`)
	matches := re.FindAllStringSubmatch(string(htmlData), -1)
	results := []string{}
	for _, m := range matches {
		results = append(results, m[0])
	}
	return results
}
