package crawler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// HTMLExtractor ...
func HTMLExtractor(link string, projectPath string) {
	// create and write files to our project directory
	// write as it downloads and not load the whole file into memory.
	// out, err := f.Create(document)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer out.Close()
	// Write the body to file
	// _, err = io.Copy(out, resp.Body)
	// fmt.Println(err)
	fmt.Println("Extracting --> ", link)
	// get the html body
	resp, err := http.Get(link)
	if err != nil {
		panic(err)
	}
	// Closure
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
