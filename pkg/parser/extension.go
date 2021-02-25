package parser

import (
	"path"
)

// URLExtension returns the extension from a given URL, or an empty string
func URLExtension(URL string) string {
	/*
		>>> https://tesla.com/main.css
		<<< ".css"

		>>> https://tesla.com/main.css?Asf341
		<<< ".css"

		>>> https://dribbble.com/css/home
		<<< ""
	*/

	// get the raw extension from the URL
	ext := path.Ext(URL)
	// fmt.Println(ext)

	// check if the extension has more than 5 characters, we need to remove any excess
	if len(ext) > 5 {
		// for every index and letter after the 0 index (ext[0] is ".", we want to keep that)
		for i, char := range ext[1:] {
			// if the unicode value of the char is not within the bounds alpha chars
			if !('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z') {
				// reinitialize the ext to only include the characters for the extension and break
				ext = ext[:i+1]
				break
			}
		}
	}

	// // return the extension
	return ext
}
