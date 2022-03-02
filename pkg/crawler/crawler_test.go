package crawler

import (
	"context"
	"os"
	"testing"

	"github.com/gocolly/colly/v2"
	"github.com/imthaghost/goclone/pkg/file"
	"github.com/imthaghost/goclone/testutils"
)

var TsUrl string

func collectAndGetFileContent(tsUrl, projectDirectory, relativeRoute string) string {
	Collector(context.Background(), tsUrl, projectDirectory, nil, "", "")
	route := projectDirectory + relativeRoute
	fileContent := file.GetFileContent(route)
	return fileContent
}

var collectorTests = map[string]func(*testing.T){
	"indexDownload": func(t *testing.T) {
		projectDirectory := file.CreateProject("test")
		collectorContent := collectAndGetFileContent(TsUrl+"/hello", projectDirectory, "/index.html")

		if collectorContent != testutils.CrawlerHelloContent {
			t.Fatalf("Expect \"%s\", but got: %s", testutils.CrawlerHelloContent, collectorContent)
		}
		os.RemoveAll(projectDirectory)
	},
	"cssDownload": func(t *testing.T) {
		projectDirectory := file.CreateProject("test")
		cssFileContent := collectAndGetFileContent(TsUrl, projectDirectory, "/css/index.css")
		if cssFileContent != testutils.CrawlerCssContent {
			t.Fatalf("Expect \"%s\", but got: %s", testutils.CrawlerCssContent, cssFileContent)
		}
		os.RemoveAll(projectDirectory)
	},
	"jsDownload": func(t *testing.T) {
		projectDirectory := file.CreateProject("test")
		jsFileContent := collectAndGetFileContent(TsUrl, projectDirectory, "/js/index.js")

		if jsFileContent != testutils.CrawlerJsContent {
			t.Fatalf("Expect \"%s\", but got: %s", testutils.CrawlerJsContent, jsFileContent)
		}
		os.RemoveAll(projectDirectory)
	},
	"imgDownload": func(t *testing.T) {
		projectDirectory := file.CreateProject("test")
		imgFileContent := collectAndGetFileContent(TsUrl, projectDirectory, "/imgs/image.png")
		if imgFileContent != testutils.CrawlerImgContent {
			t.Fatalf("Expect \"%s\", but got: %s", testutils.CrawlerImgContent, imgFileContent)
		}
		os.RemoveAll(projectDirectory)
	},
}

func TestSetUpCollector(t *testing.T) {
	testutils.SilenceStdoutInTests()
	c := colly.NewCollector(colly.Async(true))
	userAgent := "Firefox"
	setUpCollector(c, nil, nil, "http://127.0.0.1:9999", "Firefox")

	if c.UserAgent != userAgent {
		t.Fatalf("Expect %s, but got: %s", userAgent, c.UserAgent)
	}
}
func TestCollectorTests(t *testing.T) {
	testutils.SilenceStdoutInTests()
	ts := testutils.NewCrawlerTestServer()
	defer ts.Close()
	TsUrl = ts.URL
	for testName, testFuntion := range collectorTests {
		t.Run(testName, testFuntion)
	}
}
