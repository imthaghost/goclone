package crawler

// Crawl asks the neccessary crawlers for collecting links for building the web page
func Crawl(site string, projectPath string) {
	// searches for css, js, and images within a given link
	Collector(site, projectPath)
}
