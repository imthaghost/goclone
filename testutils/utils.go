package testutils

import (
	"net/http"
	"net/http/httptest"
	"os"
)

var ArrangeIndexContent = `<html>
<link rel="stylesheet" href="index.css">
<script src="index.js"></script>
<img src="image.png" alt="Red dot" />
</html>`

func NewArrangeTestServer() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(ArrangeIndexContent))
	})
	return httptest.NewServer(mux)
}

var CrawlerHelloContent = "Hello World!"
var CrawlerIndexContent = `<html>
	<link rel="stylesheet" href="index.css">
	<script src="index.js"></script>
	<img src="image.png" alt="Red dot" />
</html>`
var CrawlerCssContent = `img {
position: absolute;
left: 0px;
top: 0px;
z-index: -1;
}`
var CrawlerJsContent = `console.log('from website');`
var CrawlerImgContent = "PNG\\r\\n\\x1a\\n\\x00\\x00\\x00\\rIHDR\\x00"

func NewCrawlerTestServer() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(CrawlerHelloContent))
	})

	mux.HandleFunc("/index.css", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(CrawlerCssContent))
	})

	mux.HandleFunc("/index.js", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(CrawlerJsContent))
	})

	mux.HandleFunc("/image.png", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(CrawlerImgContent))
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(CrawlerIndexContent))
	})
	return httptest.NewServer(mux)
}

func SilenceStdoutInTests() {
	os.Stdout, _ = os.Open(os.DevNull)
}
