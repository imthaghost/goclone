package parser

import (
	"fmt"
	"testing"

	"github.com/fatih/color"
)

func TestGetDomain(t *testing.T) {
	tables := []struct {
		url      string
		expected string
	}{
		{"https://google.com", "google.com"},
	}
	for _, table := range tables {
		result := GetDomain(table.url)
		expectedresult := table.expected
		if result != expectedresult {
			t.Error()
			red := color.New(color.FgRed).SprintFunc()
			fmt.Printf("%s GetDomain Failed: %s , expected %s got %s \n", red("[-]"), table.url, expectedresult, result)

		} else {
			green := color.New(color.FgGreen).SprintFunc()
			fmt.Printf("%s GetDomain Passing: %s \n", green("[+]"), table.url)
		}
	}
}
