// Package mirror captures pages and their dependencies without flattening asset paths.
package mirror

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Options configures a sequential, bounded capture. Client is used for every request.
type Options struct {
	Client    *http.Client
	UserAgent string
	Delay     time.Duration
	MaxPages  int
	Pages     []string
}

// Resource records successful files and failures so incomplete captures are visible.
type Resource struct {
	URL         string `json:"url"`
	Path        string `json:"path,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	Error       string `json:"error,omitempty"`
}

type capture struct {
	ctx         context.Context
	start       *url.URL
	dir         string
	opts        Options
	resources   map[string]*Resource
	pages       []*url.URL
	queued      map[string]bool
	lastRequest time.Time
	bases       map[string]*url.URL
}

// Run mirrors the entry page, linked same-origin pages, and HTML/CSS dependencies.
// The output is intended to be served over HTTP from dir, not opened as file://.
func Run(ctx context.Context, site, dir string, opts Options) error {
	start, err := url.Parse(site)
	if err != nil || start == nil || (start.Scheme != "http" && start.Scheme != "https") || start.Host == "" {
		return fmt.Errorf("invalid HTTP URL: %q", site)
	}
	start.Fragment = ""
	if start.Path == "" {
		start.Path = "/"
	}
	if opts.Client == nil {
		opts.Client = &http.Client{Timeout: 30 * time.Second}
	}
	if opts.MaxPages <= 0 {
		opts.MaxPages = 100
	}
	if opts.Delay < 0 {
		return fmt.Errorf("delay must not be negative")
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	c := &capture{ctx: ctx, start: start, dir: dir, opts: opts, resources: map[string]*Resource{}, queued: map[string]bool{start.String(): true}, pages: []*url.URL{start}, bases: map[string]*url.URL{}}
	for _, raw := range opts.Pages {
		u := resolve(start, raw)
		if u == nil || u.Scheme != start.Scheme || u.Host != start.Host {
			return fmt.Errorf("additional page must have the same origin: %q", raw)
		}
		if !c.queued[u.String()] && len(c.pages) < opts.MaxPages {
			c.queued[u.String()] = true
			c.pages = append(c.pages, u)
		}
	}
	for i := 0; i < len(c.pages); i++ {
		if _, err := c.fetch(c.pages[i], true); err != nil && i == 0 {
			_ = c.manifest()
			return err
		}
		if err := ctx.Err(); err != nil {
			_ = c.manifest()
			return err
		}
	}
	if err := c.rewrite(); err != nil {
		return err
	}
	if err := c.manifest(); err != nil {
		return err
	}
	failures := 0
	for _, r := range c.resources {
		if r.Error != "" {
			failures++
		}
	}
	if failures > 0 {
		return fmt.Errorf("capture incomplete: %d failed resources (see mirror-manifest.json)", failures)
	}
	return nil
}

func (c *capture) fetch(u *url.URL, page bool) (string, error) {
	key := u.String()
	if r, ok := c.resources[key]; ok {
		if r.Error != "" {
			return "", fmt.Errorf("%s", r.Error)
		}
		return r.Path, nil
	}
	r := &Resource{URL: key}
	c.resources[key] = r
	fail := func(err error) (string, error) { r.Error = err.Error(); return "", err }
	if wait := c.opts.Delay - time.Since(c.lastRequest); wait > 0 {
		timer := time.NewTimer(wait)
		defer timer.Stop()
		select {
		case <-timer.C:
		case <-c.ctx.Done():
			return fail(c.ctx.Err())
		}
	}
	req, err := http.NewRequestWithContext(c.ctx, http.MethodGet, key, nil)
	if err != nil {
		return fail(err)
	}
	if c.opts.UserAgent != "" {
		req.Header.Set("User-Agent", c.opts.UserAgent)
	}
	c.lastRequest = time.Now()
	resp, err := c.opts.Client.Do(req)
	c.lastRequest = time.Now()
	if err != nil {
		return fail(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fail(fmt.Errorf("GET %s: HTTP %d", key, resp.StatusCode))
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 64<<20+1))
	if err != nil {
		return fail(err)
	}
	if len(body) > 64<<20 {
		return fail(fmt.Errorf("GET %s: exceeds 64 MiB resource limit", key))
	}
	r.ContentType = resp.Header.Get("Content-Type")
	if page && !strings.Contains(strings.ToLower(r.ContentType), "text/html") {
		return fail(fmt.Errorf("GET %s: expected HTML, got %s", key, r.ContentType))
	}
	r.Path = localPath(u, c.start, page, r.ContentType)
	if err := c.write(r.Path, body); err != nil {
		return fail(err)
	}
	fmt.Printf("Saved %s -> %s\n", key, r.Path)
	// Resolve references against the final response URL (redirects may change directories).
	base := resp.Request.URL
	c.bases[key] = base
	if page {
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
		if err != nil {
			return fail(err)
		}
		if href, ok := doc.Find("base[href]").First().Attr("href"); ok {
			if b := resolve(base, href); b != nil {
				base = b
			}
		}
		c.discoverHTML(doc, base)
	} else if isCSS(u, r.ContentType) {
		cssReferences(string(body), func(raw string) string { c.asset(base, raw); return raw })
	}
	return r.Path, nil
}

func (c *capture) write(name string, body []byte) error {
	file := filepath.Join(c.dir, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		return err
	}
	return os.WriteFile(file, body, 0644)
}

func (c *capture) asset(base *url.URL, raw string) {
	if u := resolve(base, raw); u != nil {
		_, _ = c.fetch(u, false)
	}
}

func (c *capture) discoverHTML(doc *goquery.Document, base *url.URL) {
	visitHTMLReferences(doc, func(raw string) string { c.asset(base, raw); return raw })
	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		raw, _ := s.Attr("href")
		u := resolve(base, raw)
		if u == nil || u.Host != c.start.Host || u.Scheme != c.start.Scheme || u.RawQuery != "" {
			return
		}
		ext := strings.ToLower(path.Ext(u.Path))
		if ext != "" && ext != ".html" && ext != ".htm" && ext != ".php" {
			c.asset(base, raw)
			return
		}
		for _, p := range []string{"/wp-admin", "/wp-json", "/wp-login.php", "/xmlrpc.php"} {
			if strings.HasPrefix(u.Path, p) {
				return
			}
		}
		if !c.queued[u.String()] && len(c.pages) < c.opts.MaxPages {
			c.queued[u.String()] = true
			c.pages = append(c.pages, u)
		}
	})
}

func resolve(base *url.URL, raw string) *url.URL {
	raw = strings.TrimSpace(raw)
	if raw == "" || strings.HasPrefix(raw, "#") {
		return nil
	}
	u, err := url.Parse(raw)
	if err != nil {
		return nil
	}
	u = base.ResolveReference(u)
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil
	}
	u.Fragment = ""
	if u.Path == "" {
		u.Path = "/"
	}
	return u
}

func localPath(u, start *url.URL, page bool, contentType string) string {
	// EscapedPath keeps encoded separators from becoming filesystem traversal.
	p := strings.TrimPrefix(path.Clean("/"+u.EscapedPath()), "/")
	if p == "" || p == "." {
		p = "index.html"
	} else if page && (strings.HasSuffix(u.Path, "/") || path.Ext(p) == "") {
		p = path.Join(p, "index.html")
	}
	if !page && path.Ext(p) == "" && isCSS(u, contentType) {
		p += ".css"
	}
	if u.RawQuery != "" {
		hash := fmt.Sprintf("%x", sha256.Sum256([]byte(u.RawQuery)))[:12]
		ext := path.Ext(p)
		p = strings.TrimSuffix(p, ext) + "__q_" + hash + ext
	}
	if u.Host != start.Host || u.Scheme != start.Scheme {
		hash := fmt.Sprintf("%x", sha256.Sum256([]byte(u.Scheme+"://"+u.Host)))[:12]
		p = path.Join("_external", hash, p)
	}
	return p
}

func isCSS(u *url.URL, contentType string) bool {
	return strings.HasSuffix(strings.ToLower(u.Path), ".css") || strings.Contains(strings.ToLower(contentType), "text/css")
}

func (c *capture) manifest() error {
	rows := make([]*Resource, 0, len(c.resources))
	for _, r := range c.resources {
		rows = append(rows, r)
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].URL < rows[j].URL })
	body, err := json.MarshalIndent(rows, "", "  ")
	if err != nil {
		return err
	}
	return c.write("mirror-manifest.json", append(body, '\n'))
}
