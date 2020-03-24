package flags

import (
	"flag"
)

// ParseFlags ...
func ParseFlags() (bool, bool, bool) {
	usage := `

	Usage:
	goclone
	goclone <url>
	goclone --help
	goclone --version
	goclone --v

Options:
	<url>  Optional url argument.
	--help  Show help screen.
	--version  Show version.
	--v		Verbose output`
	// help flag
	var help bool
	// verbose boolean flag
	var login bool
	// serve
	var serve bool
	// define help flag
	flag.BoolVar(&help, "help", false, usage)
	// define verbose flag
	flag.BoolVar(&login, "login", false, "Use this flag if you want to pass credentials to the site")
	// define serve
	flag.BoolVar(&serve, "serve", false, "Use this flag to host this site on an echo server")
	// parse the flags
	flag.Parse()
	return help, login, serve
}
