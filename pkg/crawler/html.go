package crawler

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
)

// HTMLExtractor ...
func HTMLExtractor(link string, projectPath string, referer string, userAgent string, cookieJar *cookiejar.Jar) {
	fmt.Println("Extracting --> ", link)
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	// get the html body
	client := &http.Client{Jar: cookieJar}
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return
	}
	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	}
	if referer != "" {
		req.Header.Set("Referer", referer) // 例如设置 Authorization 头
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	// Close the body once everything else is compled
	defer resp.Body.Close()
	// get the project name and path we use the path to
	f, err := os.OpenFile(projectPath+"/"+"index.html", os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	htmlData, err := io.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}
	f.Write(htmlData)

}
