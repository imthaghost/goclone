package parser

import "path"

// URLExtension returns the extension from a given URL, or an empty string
func URLExtension(URL string) string {
	/*
		>>> https://tesla.com/main.css
		<<< ".css"

		>>> https://dribbble.com/css/home
		<<< ""
	*/
	return path.Ext(URL)
}
