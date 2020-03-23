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

// Extractor visits a link dtermines if its a page or sublink downloads
// the contents to a correct directory in project folder
func Extractor(link string, projectPath string) {
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
	// file base
	base := path.Base(link)
	// file extension
	extension := filepath.Ext(base)
	// css extension
	if strings.Contains(extension, ".css") {
		var name = base[0 : len(base)-len(extension)]
		document := name + ".css"
		// get the project name and path we use the path to
		f, err := os.OpenFile(projectPath+"/"+"css/"+document, os.O_RDWR|os.O_CREATE, 0777)
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
	// js extension
	if strings.Contains(extension, ".js") {
		var name = base[0 : len(base)-len(extension)]
		document := name + ".js"

		f, err := os.OpenFile(projectPath+"/"+"js/"+document, os.O_RDWR|os.O_CREATE, 0777)
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
	// jpg extension
	if strings.Contains(extension, ".jpg") {
		var name = base[0 : len(base)-len(extension)]
		document := name + ".jpg"

		f, err := os.OpenFile(projectPath+"/"+"imgs/"+document, os.O_RDWR|os.O_CREATE, 0777)
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
	// png entension
	if strings.Contains(extension, ".png") {
		var name = base[0 : len(base)-len(extension)]
		document := name + ".png"

		f, err := os.OpenFile(projectPath+"/"+"imgs/"+document, os.O_RDWR|os.O_CREATE, 0777)
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
	// gif extension
	if strings.Contains(extension, ".gif") {
		var name = base[0 : len(base)-len(extension)]
		document := name + ".gif"

		f, err := os.OpenFile(projectPath+"/"+"imgs/"+document, os.O_RDWR|os.O_CREATE, 0777)
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
	// jpeg extension
	if strings.Contains(extension, ".jpeg") {
		var name = base[0 : len(base)-len(extension)]
		document := name + ".jpeg"

		f, err := os.OpenFile(projectPath+"/"+"imgs/"+document, os.O_RDWR|os.O_CREATE, 0777)
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
	// svg extension
	if strings.Contains(extension, ".svg") {
		var name = base[0 : len(base)-len(extension)]
		document := name + ".svg"

		f, err := os.OpenFile(projectPath+"/"+"imgs/"+document, os.O_RDWR|os.O_CREATE, 0777)
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
}
