package parser

import (
	"net/url"

	strutil "github.com/torden/go-strutil"
)

// ValidateURL checks for a valid url
func ValidateURL(url string) bool {
	/*
		>>> https://google.com
		<<< true

		>>> google.com
		<<< false
	*/
	if !strutil.NewStringValidator().IsValidURL(url) {
		return false
	}

	return true
}

// ValidateDomain checks for a valid domain
func ValidateDomain(domain string) bool {
	/*
		>>> google.com
		<<< true

		>>> google
		<<< false
	*/
	if !strutil.NewStringValidator().IsValidDomain(domain) {
		return false
	}

	return true
}

// CreateURL will take in a valid domain and return the URL
func CreateURL(domain string) string {
	/*
		>>> google.com - Valid Domain
		<<< https://google.com - Returned URL
	*/

	// concate nate https:// and the valid domain and return the now valid url
	return "https://" + domain
}

// GetDomain takes in a valid URL and returns the domain of the url
func GetDomain(validurl string) string {
	/*
		>>> https://google.com - Valid URL
		<<< google.com - Hostname
	*/

	// parse the url
	u, err := url.Parse(validurl)

	if err != nil {
		panic(err)
	}

	// grab the hostname from the string
	hostname := u.Hostname()

	// hostname
	return hostname
}
