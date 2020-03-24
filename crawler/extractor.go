package crawler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Extractor visits a link determines if its a page or sublink downloads
// the contents to a correct directory in project folder
func Extractor(link string, projectPath string) {
	fmt.Println("Extracting --> ", link)

	// get the html body
	resp, err := http.Get(link)
	if err != nil {
		panic(err)
	}

	// Closure
	defer resp.Body.Close()
	// file base
	base := path.Base(link)
	// file extension
	extension := filepath.Ext(base)

	// I wish we could use a switch statement..
	if strings.Contains(extension, ".css") {
		writeFileToPath(projectPath, base, extension, ".css", "css", resp)
	} else if strings.Contains(extension, ".js") {
		writeFileToPath(projectPath, base, extension, ".js", "js", resp)
	} else if strings.Contains(extension, ".jpg") {
		writeFileToPath(projectPath, base, extension, ".jpg", "imgs", resp)
	} else if strings.Contains(extension, ".jpeg") {
		writeFileToPath(projectPath, base, extension, ".jpeg", "imgs", resp)
	} else if strings.Contains(extension, ".gif") {
		writeFileToPath(projectPath, base, extension, ".gif", "imgs", resp)
	} else if strings.Contains(extension, ".svg") {
		writeFileToPath(projectPath, base, extension, ".svg", "imgs", resp)
	} else if strings.Contains(extension, ".png") {
		writeFileToPath(projectPath, base, extension, ".png", "imgs", resp)
	}
}

func writeFileToPath(projectPath, base, oldFileExt, newFileExt, fileDir string, resp *http.Response) {
	var name = base[0 : len(base)-len(oldFileExt)]
	document := name + newFileExt

	// get the project name and path we use the path to
	f, err := os.OpenFile(projectPath+"/"+fileDir+"/"+document, os.O_RDWR|os.O_CREATE, 0777)
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
