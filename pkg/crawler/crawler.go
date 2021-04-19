package crawler

import (
	"context"

	"github.com/gocolly/colly"
)

// Crawl asks the necessary crawlers for collecting links for building the web page
func Crawl(ctx context.Context, site string, projectPath string, collyOpts ...func(*colly.Collector)) error {
	// searches for css, js, and images within a given link
	return Collector(ctx, site, projectPath, collyOpts...)
}
