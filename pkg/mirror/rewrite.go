package mirror

import (
	"bytes"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	xhtml "golang.org/x/net/html"
)

var cssURL = regexp.MustCompile(`(?i)url\(\s*(?:"([^"]*)"|'([^']*)'|([^\s)'";][^)]*))\s*\)`)
var cssImport = regexp.MustCompile(`(?i)(@import\s+)(?:"([^"]+)"|'([^']+)')`)

func cssReferences(css string, replace func(string) string) string {
	css = cssURL.ReplaceAllStringFunc(css, func(match string) string {
		m := cssURL.FindStringSubmatch(match)
		raw := m[1]
		if raw == "" {
			raw = m[2]
		}
		if raw == "" {
			raw = strings.TrimSpace(m[3])
		}
		changed := replace(raw)
		if changed == raw {
			return match
		}
		return `url("` + changed + `")`
	})
	return cssImport.ReplaceAllStringFunc(css, func(match string) string {
		m := cssImport.FindStringSubmatch(match)
		raw := m[2]
		if raw == "" {
			raw = m[3]
		}
		return m[1] + `"` + replace(raw) + `"`
	})
}

func srcsetReferences(value string, replace func(string) string) string {
	// Data URLs contain commas; leave them intact rather than inventing fetches.
	if strings.Contains(value, "data:") {
		return value
	}
	parts := strings.Split(value, ",")
	for i, part := range parts {
		fields := strings.Fields(part)
		if len(fields) > 0 {
			fields[0] = replace(fields[0])
			parts[i] = strings.Join(fields, " ")
		}
	}
	return strings.Join(parts, ", ")
}

func visitHTMLReferences(doc *goquery.Document, replace func(string) string) {
	doc.Find("*").Each(func(_ int, s *goquery.Selection) {
		for _, attr := range []string{"src", "poster", "data-src", "data-lazy-src", "data-splide-lazy", "data-bg", "data-background-image", "data-full-image", "data-light-image"} {
			if raw, ok := s.Attr(attr); ok {
				// Embedded documents are independent applications, not static assets.
				if s.Is("iframe") {
					continue
				}
				s.SetAttr(attr, replace(raw))
			}
		}
		for _, attr := range []string{"srcset", "data-srcset", "data-lazy-srcset", "data-splide-lazy-srcset"} {
			if raw, ok := s.Attr(attr); ok {
				s.SetAttr(attr, srcsetReferences(raw, replace))
			}
		}
		if raw, ok := s.Attr("style"); ok {
			s.SetAttr("style", cssReferences(raw, replace))
		}
	})
	doc.Find("link[href]").Each(func(_ int, s *goquery.Selection) {
		rel, _ := s.Attr("rel")
		as, _ := s.Attr("as")
		if strings.Contains(strings.ToLower(rel), "stylesheet") || strings.Contains(strings.ToLower(rel), "icon") || as == "style" || as == "image" || as == "font" || as == "script" {
			raw, _ := s.Attr("href")
			s.SetAttr("href", replace(raw))
		}
	})
	doc.Find("style").Each(func(_ int, s *goquery.Selection) {
		// Browsers do not decode HTML entities inside raw-text style elements.
		css := cssReferences(s.Text(), replace)
		for _, node := range s.Nodes {
			for node.FirstChild != nil {
				node.RemoveChild(node.FirstChild)
			}
			node.AppendChild(&xhtml.Node{Type: xhtml.TextNode, Data: css})
		}
	})
}

func (c *capture) rewrite() error {
	for _, r := range c.resources {
		if r.Error != "" {
			continue
		}
		u, _ := url.Parse(r.URL)
		page := c.queued[r.URL]
		css := isCSS(u, r.ContentType)
		js := strings.Contains(r.ContentType, "javascript") || strings.HasSuffix(u.Path, ".js")
		if !page && !css && !js {
			continue
		}
		file := filepath.Join(c.dir, filepath.FromSlash(r.Path))
		body, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		base := c.bases[r.URL]
		if base == nil {
			base = u
		}
		replace := func(raw string) string {
			resolved := resolve(base, raw)
			if resolved == nil {
				return raw
			}
			if target, ok := c.resources[resolved.String()]; ok && target.Error == "" && target.Path != "" {
				fragment := ""
				if original, err := url.Parse(raw); err == nil && original.Fragment != "" {
					fragment = "#" + original.EscapedFragment()
				}
				return (&url.URL{Path: "/" + target.Path}).EscapedPath() + fragment
			}
			return raw
		}
		if page {
			doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
			if err != nil {
				return err
			}
			if href, ok := doc.Find("base[href]").First().Attr("href"); ok {
				if b := resolve(base, href); b != nil {
					base = b
				}
			}
			visitHTMLReferences(doc, replace)
			doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) { raw, _ := s.Attr("href"); s.SetAttr("href", replace(raw)) })
			doc.Find("base").Remove()
			output, err := doc.Html()
			if err != nil {
				return err
			}
			body = []byte(output)
		} else if css {
			body = []byte(cssReferences(string(body), replace))
		}
		// WordPress plugins embed their asset base in inline JSON and JS strings.
		origin := c.start.Scheme + "://" + c.start.Host
		body = bytes.ReplaceAll(body, []byte(origin+"/"), []byte("/"))
		body = bytes.ReplaceAll(body, []byte(strings.ReplaceAll(origin+"/", "/", `\/`)), []byte(`\/`))
		if err := c.write(r.Path, body); err != nil {
			return err
		}
	}
	return nil
}
