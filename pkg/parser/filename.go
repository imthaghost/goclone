package parser

import "path"

// URLFilename returns the file name from a given url
func URLFilename(filename string) string {
	/*
		>>> https://tesla.com/main.css
		<<< main.css

		>>> https://dribbble.com/css/home.css
		<<< home.css
	*/
	return path.Base(filename)
}

// PathFilename returns the file name from a given path
func PathFilename(givenPath string) string {
	/*
		>>> /css/main.css
		<<< main.css

		>>> /js/googleanalytics.js
		<<< googleanalytics.js
	*/
	return path.Base(givenPath)
}
