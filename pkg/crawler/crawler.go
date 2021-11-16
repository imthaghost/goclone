package crawler

import (
	"context"
	"net/http/cookiejar"
)

// Crawl asks the necessary crawlers for collecting links for building the web page
func Crawl(ctx context.Context, site string, projectPath string, cookieJar *cookiejar.Jar, proxyString string, userAgent string) error {
	// searches for css, js, and images within a given link
	return Collector(ctx, site, projectPath, cookieJar, proxyString, userAgent)
}
