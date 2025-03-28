package html

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/goclone-dev/goclone/pkg/crawler"
	"github.com/goclone-dev/goclone/pkg/file"
	"github.com/goclone-dev/goclone/testutils"
)

// TestArrange verifies that the LinkRestructure function correctly reorganizes the paths
// of resources (CSS, JS and images) in the HTML file, ensuring that:
// 1. Paths are correctly updated to their new locations
// 2. Files exist in their expected locations
// 3. Original element attributes are preserved
func TestArrange(t *testing.T) {
	testutils.SilenceStdoutInTests()
	ts := testutils.NewArrangeTestServer()
	defer ts.Close()

	// Initial setup
	projectDirectory := file.CreateProject("test")
	defer os.RemoveAll(projectDirectory)

	// Run crawler and restructuring
	crawler.Collector(context.Background(), ts.URL, projectDirectory, nil, "", "")

	if err := LinkRestructure(projectDirectory); err != nil {
		t.Fatalf("Error during restructuring: %v", err)
	}

	// Verify that index.html exists
	if !file.Exists(projectDirectory + "/index.html") {
		t.Fatal("index.html file should exist")
	}

	// Get and verify file content
	indexFileContent := file.GetFileContent(projectDirectory + "/index.html")
	if indexFileContent == testutils.ArrangeIndexContent {
		t.Fatalf("Expected restructured HTML, not original: %s", testutils.ArrangeIndexContent)
	}

	// Verify that files exist in expected locations
	expectedFiles := []string{
		"/css/index.css",
		"/js/index.js",
		"/imgs/image.png",
	}

	for _, expectedFile := range expectedFiles {
		if !file.Exists(projectDirectory + expectedFile) {
			t.Fatalf("File %s should exist", expectedFile)
		}
	}

	// Verify paths in HTML
	expectedPaths := []string{
		"css/index.css",
		"js/index.js",
		"imgs/image.png",
	}

	for _, path := range expectedPaths {
		if !strings.Contains(indexFileContent, path) {
			t.Fatalf("Expected to find path %s in HTML", path)
		}
	}

	// Verify that original attributes are preserved
	if !strings.Contains(indexFileContent, `alt="Red dot"`) {
		t.Fatal("Expected to preserve alt attribute in image")
	}
}
