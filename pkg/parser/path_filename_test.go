package parser

import (
	"fmt"
	"testing"

	"github.com/fatih/color"
)

func TestPathFilename(t *testing.T) {
	tables := []struct {
		path     string
		expected string
	}{
		{"/css/main.css", "main.css"},
		{"/js/googleanalytics.js", "googleanalytics.js"},
		{"/jquery/min.js", "min.js"},
	}
	for _, table := range tables {
		result := PathFilename(table.path)
		expectedresult := table.expected
		if result != expectedresult {
			t.Error()
			red := color.New(color.FgRed).SprintFunc()
			fmt.Printf("%s PathFilename Failed: %s , expected %s got %s \n", red("[-]"), table.path, expectedresult, result)

		} else {
			green := color.New(color.FgGreen).SprintFunc()
			fmt.Printf("%s PathFilename Passing: %s \n", green("[+]"), table.path)
		}
	}
}
