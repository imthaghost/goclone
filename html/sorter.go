package html

import (
	"bytes"

	"golang.org/x/net/html"
)

// LinkRestructure grabs all html files in project directory
// reorganizes each file with local links (css js images)
func LinkRestructure(projectDir string) {
	// arrange css
	// arrangeCSS(projectDir)
	// // arrange js
	// arrangeJS(projectDir)
	// // arrange imgs
	// arrangeImgs(projectDir)

	arrange(projectDir)
}

func collectText(n *html.Node, buf *bytes.Buffer) {
	if n.Type == html.TextNode {
		buf.WriteString(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		collectText(c, buf)
	}
}
