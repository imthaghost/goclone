package mirror

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

// Regression for the external-asset case raised in PR #30: downloading a CDN
// script is insufficient if the generated HTML still points at the CDN.
func TestExternalAssetsWorkOffline(t *testing.T) {
	cdn := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/scripts/vendor.js":
			w.Header().Set("Content-Type", "application/javascript")
			fmt.Fprint(w, "window.externalAsset = true;")
		case "/styles/main.css":
			w.Header().Set("Content-Type", "text/css")
			fmt.Fprint(w, `@font-face{src:url('../fonts/test.woff2')}`)
		case "/fonts/test.woff2":
			fmt.Fprint(w, "font fixture")
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(cdn.Close)
	origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/scripts/vendor.js" {
			w.Header().Set("Content-Type", "application/javascript")
			fmt.Fprint(w, "window.localAsset = true;")
			return
		}
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `<script src="/scripts/vendor.js"></script><script src="%s/scripts/vendor.js"></script><link rel="stylesheet" href="%s/styles/main.css">`, cdn.URL, cdn.URL)
	}))
	t.Cleanup(origin.Close)
	dir := t.TempDir()
	if err := Run(context.Background(), origin.URL, dir, Options{MaxPages: 1}); err != nil {
		t.Fatal(err)
	}
	// Both remote origins are gone before the captured files are served.
	origin.Close()
	cdn.Close()
	local := httptest.NewServer(http.FileServer(http.Dir(dir)))
	defer local.Close()
	get := func(ref string) string {
		t.Helper()
		if !strings.HasPrefix(ref, "/") || strings.HasPrefix(ref, "//") {
			t.Fatalf("capture still needs a remote asset: %q", ref)
		}
		resp, err := local.Client().Get(local.URL + ref)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("local resource %s: status=%d err=%v", ref, resp.StatusCode, err)
		}
		return string(body)
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(get("/index.html")))
	if err != nil {
		t.Fatal(err)
	}
	scripts := doc.Find("script[src]")
	if scripts.Length() != 2 {
		t.Fatalf("expected both scripts, got %d", scripts.Length())
	}
	for i, expected := range []string{"window.localAsset = true;", "window.externalAsset = true;"} {
		ref, _ := scripts.Eq(i).Attr("src")
		if body := get(ref); body != expected {
			t.Fatalf("same-basename script overwritten: got %q, want %q", body, expected)
		}
	}
	stylesheet, _ := doc.Find("link[rel=stylesheet]").Attr("href")
	css := get(stylesheet)
	// Follow the font reference actually emitted in the served stylesheet.
	_, after, ok := strings.Cut(css, `url("`)
	if !ok {
		t.Fatalf("missing local font URL: %s", css)
	}
	font, _, ok := strings.Cut(after, `")`)
	if !ok || get(font) != "font fixture" {
		t.Fatalf("external CSS font is not available offline: %s", css)
	}
}
