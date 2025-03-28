package parser

import (
	"fmt"
	"testing"

	"github.com/fatih/color"
)

func TestURLExtension(t *testing.T) {
	tables := []struct {
		url      string
		expected string
	}{
		// Basic cases
		{"https://tesla.com/main.css", ".css"},
		{"https://tesla.com/main.js", ".js"},
		{"https://tesla.com/image.jpg", ".jpg"},
		{"https://tesla.com/image.png", ".png"},
		{"https://tesla.com/image.svg", ".svg"},

		// Cases with query parameters
		{"https://tesla.com/main.css?Asf341", ".css"},
		{"https://tesla.com/main.css?ver=1.0", ".css"},
		{"https://tesla.com/main.css?ver=1.0&type=min", ".css"},

		// Cases with fragments
		{"https://tesla.com/main.css#section", ".css"},
		{"https://tesla.com/main.css#top", ".css"},

		// Combined cases
		{"https://tesla.com/main.css?ver=1.0#section", ".css"},
		{"https://tesla.com/main.min.js?ver=1.0#section", ".js"},

		// Cases without extension
		{"https://dribbble.com/css/home", ""},
		{"https://example.com/path/to/file", ""},
		{"https://example.com/path/to/file/", ""},

		// Cases with multiple dots
		{"https://example.com/file.min.js", ".js"},
		{"https://example.com/file.min.css", ".css"},

		// Cases with long paths
		{"https://example.com/path/to/file.min.js", ".js"},
		{"https://example.com/path/to/deep/file.min.css", ".css"},

		// Cases with special characters
		{"https://example.com/file-name.min.js", ".js"},
		{"https://example.com/file_name.min.css", ".css"},

		// Cases with malformed URLs
		{"", ""},
		{"not-a-url", ""},
		{"http://", ""},
	}

	for _, table := range tables {
		result := URLExtension(table.url)
		expectedresult := table.expected
		if result != expectedresult {
			t.Error()
			red := color.New(color.FgRed).SprintFunc()
			fmt.Printf("%s URLExtension Failed: %s , expected %s got %s \n", red("[-]"), table.url, expectedresult, result)
		} else {
			green := color.New(color.FgGreen).SprintFunc()
			fmt.Printf("%s URLExtension Passing: %s \n", green("[+]"), table.url)
		}
	}
}
