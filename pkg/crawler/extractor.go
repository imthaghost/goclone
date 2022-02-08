package crawler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

// file extension map for directing files to their proper directory in O(1) time
var (
	supportedExtension = map[string]struct{}{
		".css": {},
		".js":   {},
		".jpg":  {},
		".jpeg": {},
		".gif":  {},
		".png":  {},
		".svg":  {},
		".webp":  {},
	}
)

// Extractor visits a link determines if its a page or sublink
// downloads the contents to a correct directory in project folder
// TODO add functionality for determining if page or sublink
func Extractor(link string, projectPath string) {
	fmt.Println("Extracting --> ", link)

	// get the html body
	resp, err := http.Get(link)
	if err != nil {
		panic(err)
	}

	// Closure
	defer resp.Body.Close()

	parsedURL, err := url.Parse(link)
	if err != nil {
		panic(err)
	}

	// If extension is supported then writeFileToPath
	if _, ok := supportedExtension[path.Ext(link)]; ok {
		writeFileToPath(projectPath, parsedURL.Path, resp)
	}
}


func writeFileToPath(projectPath, path string, resp *http.Response) {

	pathToWrite := filepath.Join(projectPath, path)
	// mkdir -p
	if err := os.MkdirAll(filepath.Dir(pathToWrite), os.ModePerm); err != nil {
		panic(err)
	}

	f, err := os.OpenFile(pathToWrite, os.O_RDWR|os.O_CREATE, 0777)
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
