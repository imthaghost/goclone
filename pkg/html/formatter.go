package html

import (
	"io/ioutil"

	"github.com/yosssi/gohtml"
)

// FormatHTML will formart any given string of HTML
func FormatHTML(filePath string) {
	// TODO: Implement
	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	data := gohtml.Format(string(dat))
	b := []byte(data)
	error := ioutil.WriteFile(filePath, b, 0777)
	// handle this error
	if error != nil {
		// print it out
		panic(error)
	}
}


