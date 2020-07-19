package parser

import (
	"fmt"
	"testing"

	"github.com/fatih/color"
)

func TestCreateURL(t *testing.T) {
	tables := []struct {
		domain   string
		expected string
	}{
		{"google.com", "https://google.com"},
		{"github.com", "https://github.com"},
		{"tesla.com", "https://tesla.com"},
		{"amazon.com", "https://amazon.com"},
	}
	for _, table := range tables {
		result := CreateURL(table.domain)
		expectedresult := table.expected
		if result != expectedresult {
			t.Error()
			red := color.New(color.FgRed).SprintFunc()
			fmt.Printf("%s CreateURL Failed: %s , expected %s got %s \n", red("[-]"), table.domain, expectedresult, result)

		} else {
			green := color.New(color.FgGreen).SprintFunc()
			fmt.Printf("%s CreateURL Passing: %s \n", green("[+]"), table.domain)
		}
	}
}
