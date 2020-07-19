package tests

import (
	"fmt"
	"testing"

	"github.com/fatih/color"
	"github.com/imthaghost/goclone/parser"
)

func TestPathFilename(t *testing.T) {
	tables := []struct {
		path     string
		expected string
	}{
		{"/css/main.css", "main.css"},
		{"/js/googleanalytics.js", "googleanalytics.js"},
	}
	for _, table := range tables {
		result := parser.PathFilename(table.path)
		expectedresult := table.expected
		if result != expectedresult {
			t.Error()
			red := color.New(color.FgRed).SprintFunc()
			fmt.Printf("%s URLFilename Failed: %s , expected %s got %s \n", red("[-]"), table.path, expectedresult, result)

		} else {
			green := color.New(color.FgGreen).SprintFunc()
			fmt.Printf("%s Passing: %s \n", green("[+]"), table.path)
		}
	}
}
