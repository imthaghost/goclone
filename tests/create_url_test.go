package tests

import (
	"fmt"
	"testing"

	"github.com/fatih/color"
	"github.com/imthaghost/goclone/parser"
)

func TestCreateURL(t *testing.T) {
	tables := []struct {
		domain   string
		expected string
	}{
		{"google.com", "https://google.com"},
		{"github.com", "https://github.com"},
	}
	for _, table := range tables {
		result := parser.CreateURL(table.domain)
		expectedresult := table.expected
		if result != expectedresult {
			t.Error()
			red := color.New(color.FgRed).SprintFunc()
			fmt.Printf("%s URLFilename Failed: %s , expected %s got %s \n", red("[-]"), table.domain, expectedresult, result)

		} else {
			green := color.New(color.FgGreen).SprintFunc()
			fmt.Printf("%s Passing: %s \n", green("[+]"), table.domain)
		}
	}
}
