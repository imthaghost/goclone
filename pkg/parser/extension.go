package parser

import (
	"net/url"
	"path"
)

// URLExtension returns the extension from a given URL, or an empty string
func URLExtension(URL string) string {
	// Parse the URL
	u, err := url.Parse(URL)
	if err != nil {
		return ""
	}

	// Get the path and extract the extension
	return path.Ext(u.Path)
}
