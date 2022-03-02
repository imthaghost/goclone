package html

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/imthaghost/goclone/pkg/crawler"
	"github.com/imthaghost/goclone/pkg/file"
	"github.com/imthaghost/goclone/testutils"
)

func TestArrange(t *testing.T) {
	testutils.SilenceStdoutInTests()
	ts := testutils.NewArrangeTestServer()
	defer ts.Close()
	projectDirectory := file.CreateProject("test")
	crawler.Collector(context.Background(), ts.URL, projectDirectory, nil, "", "")
	LinkRestructure(projectDirectory)
	indexFileContent := file.GetFileContent(projectDirectory + "/index.html")
	if indexFileContent == testutils.ArrangeIndexContent {
		t.Fatalf("Expect restructure html, no orginial: %s", testutils.ArrangeIndexContent)
	}
	if !strings.Contains(indexFileContent, "css/index.css") {
		fmt.Println(indexFileContent)
		t.Fatalf("Expect css route in html")
	}
	if !strings.Contains(indexFileContent, "js/index.js") {
		t.Fatalf("Expect js route in html")
	}
	if !strings.Contains(indexFileContent, "imgs/image.png") {
		t.Fatalf("Expect imgs route in html")
	}
	os.RemoveAll(projectDirectory)
}
