package crawler

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
)

// HTMLExtractor downloads the HTML content from a URL and saves it to index.html
func HTMLExtractor(link string, projectPath string) error {
	fmt.Println("Extracting HTML from --> ", link)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	// get the html body
	resp, err := http.Get(link)
	if err != nil {
		return fmt.Errorf("failed to GET HTML: %v", err)
	}

	// Close the body once everything else is compled
	defer resp.Body.Close()

	// Read the HTML content
	htmlData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read HTML body: %v", err)
	}

	// get the project name and path we use the path to
	f, err := os.OpenFile(projectPath+"/"+"index.html", os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return fmt.Errorf("failed to open index.html: %v", err)
	}
	defer f.Close()

	_, err = f.Write(htmlData)
	if err != nil {
		return fmt.Errorf("failed to write HTML content: %v", err)
	}

	return nil
}
