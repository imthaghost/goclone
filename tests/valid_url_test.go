package tests

import (
	"fmt"
	"testing"

	"github.com/fatih/color"
	"github.com/imthaghost/goclone/parser"
)

func TestValidURL(t *testing.T) {
	tables := []struct {
		url      string
		expected bool
	}{
		{"https://google.com", true},
		{"https://tesla.com", true},
		{"google.com", false},
		{"g234hj3k242", false},
		{"ws://google.com", false},
		{"https://amazon.com", true},
		{"https://medium.com", true},
		{"https://www.codemio.com/2018/03/http2-streaming-golang.html", true},
		{"https://github.com/sentriz/gonic", true},
		{"https://echo.labstack.com/cookbook/auto-tls", true},
	}
	for _, table := range tables {
		result := parser.ValidateURL(table.url)
		expectedresult := table.expected
		if result != expectedresult {
			t.Error()
			red := color.New(color.FgRed).SprintFunc()
			fmt.Printf("%s URLFilename Failed: %s , expected %t got %t \n", red("[-]"), table.url, expectedresult, result)

		} else {
			green := color.New(color.FgGreen).SprintFunc()
			fmt.Printf("%s Passing: %s \n", green("[+]"), table.url)
		}
	}
}
