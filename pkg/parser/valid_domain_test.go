package parser

import (
	"fmt"
	"testing"

	"github.com/fatih/color"
)

func TestValidDomain(t *testing.T) {
	tables := []struct {
		domain   string
		expected bool
	}{
		{"google.com", true},
		{"tesla.com", true},
		{"g234hj3k242", false},
		{"ws://google.com", false},
		{"amazon.com", true},
		{"medium.com", true},
		{"codemio.com", true},
		{"github.com", true},
		{"echo.labstack.com", true},
	}
	for _, table := range tables {
		result := ValidateDomain(table.domain)
		expectedresult := table.expected
		if result != expectedresult {
			t.Error()
			red := color.New(color.FgRed).SprintFunc()
			fmt.Printf("%s ValidDomain Failed: %s , expected %t got %t \n", red("[-]"), table.domain, expectedresult, result)

		} else {
			green := color.New(color.FgGreen).SprintFunc()
			fmt.Printf("%s ValidDomain Passing: %s \n", green("[+]"), table.domain)
		}
	}
}
