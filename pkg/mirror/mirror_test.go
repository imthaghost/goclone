package mirror

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestWordPressDependencies(t *testing.T) {
	var mu sync.Mutex
	counts := map[string]int{}
	external := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/css" {
			w.Header().Set("Content-Type", "text/css")
			fmt.Fprint(w, `@font-face{src:url('/font.woff2')}`)
		} else {
			fmt.Fprint(w, "font")
		}
	}))
	defer external.Close()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		counts[r.URL.RequestURI()]++
		mu.Unlock()
		if r.Header.Get("User-Agent") != "mirror-test" {
			t.Errorf("user agent missing: %s", r.URL)
		}
		switch r.URL.Path {
		case "/":
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintf(w, `<!doctype html><link rel="stylesheet" href="/one/main.css?ver=1"><link rel="stylesheet" href="/two/main.css?ver=1"><link rel="stylesheet" href="%s/css?family=Inter"><img src="/hero.webp" srcset="/hero.webp 1x, /hero@2x.avif 2x"><img data-lazy-src="/lazy.webp"><img src="data:image/png;base64,AAAA"><style>.hero{background:url('/inline.webp')}body > a[data-x="yes"]:after{content:"a & b"}</style><a href="/about/">About</a><a href="/about/#team">Team</a><a href="/wp-admin/">Admin</a>`, external.URL)
		case "/about/":
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, `<img src="../hero.webp"><a href="/">Home</a>`)
		case "/one/main.css":
			w.Header().Set("Content-Type", "text/css")
			fmt.Fprint(w, `@import "nested.css";a{background:url('../bg.webp')}b{mask:url('../icon.svg#mark')}`)
		case "/one/nested.css":
			w.Header().Set("Content-Type", "text/css")
			fmt.Fprint(w, `@import url('main.css?ver=1');`)
		case "/two/main.css":
			w.Header().Set("Content-Type", "text/css")
			fmt.Fprint(w, `b{color:blue}`)
		default:
			fmt.Fprint(w, "asset")
		}
	}))
	defer srv.Close()
	for attempt := 0; attempt < 2; attempt++ {
		dir := t.TempDir()
		if err := Run(context.Background(), srv.URL+"/", dir, Options{UserAgent: "mirror-test"}); err != nil {
			t.Fatal(err)
		}
		manifest := readManifest(t, dir)
		if len(manifest) != 13 {
			t.Fatalf("got %d resources: %+v", len(manifest), manifest)
		}
		home := read(t, dir, "index.html")
		for _, expected := range []string{`href="/about/index.html"`, `href="/about/index.html#team"`, `/hero@2x.avif 2x`, `data:image/png;base64,AAAA`, `/one/main__q_`, `/two/main__q_`, `/_external/`, `body > a[data-x="yes"]:after{content:"a & b"}`, `url('/inline.webp')`} {
			if !strings.Contains(home, expected) {
				t.Errorf("missing %q in %s", expected, home)
			}
		}
		for _, r := range manifest {
			if r.Error != "" {
				t.Fatal(r.Error)
			}
			if strings.ContainsAny(r.Path, "?:\\") {
				t.Errorf("unsafe path %q", r.Path)
			}
			if strings.Contains(r.URL, "/one/main.css") {
				css := read(t, dir, r.Path)
				if !strings.Contains(css, `url("/bg.webp")`) || !strings.Contains(css, `url("/icon.svg#mark")`) || !strings.Contains(css, `@import "/one/nested.css"`) {
					t.Errorf("bad css: %s", css)
				}
			}
		}
	}
	mu.Lock()
	defer mu.Unlock()
	if counts["/hero.webp"] != 2 {
		t.Errorf("duplicate downloads or global cache leaked: %v", counts)
	}
	if counts["/wp-admin/"] != 0 {
		t.Error("crawled admin")
	}
}

func read(t *testing.T, dir, name string) string {
	t.Helper()
	b, e := os.ReadFile(filepath.Join(dir, filepath.FromSlash(name)))
	if e != nil {
		t.Fatal(e)
	}
	return string(b)
}
func readManifest(t *testing.T, dir string) []Resource {
	t.Helper()
	var rows []Resource
	if e := json.Unmarshal([]byte(read(t, dir, "mirror-manifest.json")), &rows); e != nil {
		t.Fatal(e)
	}
	return rows
}

func TestFailureAndCancellation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, `<img src="/missing.webp">`)
		} else {
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()
	dir := t.TempDir()
	if err := Run(context.Background(), srv.URL+"/", dir, Options{}); err == nil {
		t.Fatal("missing asset reported as success")
	}
	if _, err := os.Stat(filepath.Join(dir, "missing.webp")); !os.IsNotExist(err) {
		t.Fatal("saved error page as an image")
	}
	if !strings.Contains(read(t, dir, "mirror-manifest.json"), "HTTP 404") {
		t.Fatal("missing failure detail")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := Run(ctx, srv.URL+"/", t.TempDir(), Options{}); err == nil {
		t.Fatal("cancellation ignored")
	}
}

func TestBoundsRedirectBaseAndQueryVariants(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/start":
			http.Redirect(w, r, "/nested/", http.StatusFound)
		case "/nested/":
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, `<base href="assets/"><img src="photo.webp?v=1"><img src="photo.webp?v=2"><a href="/second/">two</a>`)
		case "/nested/assets/photo.webp":
			fmt.Fprint(w, r.URL.RawQuery)
		default:
			t.Errorf("unexpected page: %s", r.URL)
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()
	dir := t.TempDir()
	if err := Run(context.Background(), srv.URL+"/start", dir, Options{MaxPages: 1}); err != nil {
		t.Fatal(err)
	}
	if len(readManifest(t, dir)) != 3 {
		t.Fatal("max pages not respected")
	}
	home := read(t, dir, "start/index.html")
	if strings.Contains(home, "<base") || strings.Count(home, "/nested/assets/photo__q_") != 2 {
		t.Fatal(home)
	}
	u, _ := url.Parse(srv.URL + "/%2e%2e/%2Foutside.webp?x=one")
	start, _ := url.Parse(srv.URL)
	p := localPath(u, start, false, "")
	if strings.Contains(p, "../") || filepath.IsAbs(p) {
		t.Fatal(p)
	}
}

func TestDelayAndTLS(t *testing.T) {
	var times []time.Time
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		times = append(times, time.Now())
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<a href="/second/">two</a>`)
	}))
	defer srv.Close()
	if err := Run(context.Background(), srv.URL+"/", t.TempDir(), Options{}); err == nil {
		t.Fatal("TLS verification disabled by default")
	}
	if err := Run(context.Background(), srv.URL+"/", t.TempDir(), Options{Client: srv.Client(), Delay: 30 * time.Millisecond}); err != nil {
		t.Fatal(err)
	}
	if len(times) != 2 || times[1].Sub(times[0]) < 25*time.Millisecond {
		t.Fatalf("delay ignored: %v", times)
	}
}
