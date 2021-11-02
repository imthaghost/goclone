package crawler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"crypto/tls"
)

// HTMLExtractor ...
func HTMLExtractor(link string, projectPath string) {
	fmt.Println("Extracting --> ", link)
	
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	// get the html body
	resp, err := http.Get(link)
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
	htmlData, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}
	f.Write(htmlData)

}
