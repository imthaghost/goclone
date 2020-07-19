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
		{"https://tesla.com/main.css", ".css"},
		{"https://tesla.com/main.css?Asf341", ".css"},
		{"https://dribbble.com/css/home", ""},
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
