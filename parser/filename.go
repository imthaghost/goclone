package parser

// URLFilename returns the file name from a given url
func URLFilename(filename string) string {
	/*
		>>> https://tesla.com/main.css
		<<< main.css

		>>> https://dribbble.com/css/home.css
		<<< home.css
	*/
	some := "Nothing for now"
	return some
}

// PathFilename returns the file name from a given path
func PathFilename(path string) string {
	/*
		>>> /css/main.css
		<<< main.css

		>>> /js/googleanalytics.js
		<<< googleanalytics.js
	*/
	some := "Nothing for now"
	return some
}
