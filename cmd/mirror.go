package cmd

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

// hostTransport keeps explicit certificate exceptions off the shared transport
// and re-evaluates the destination on every request, including redirects.
type hostTransport struct {
	verified  *http.Transport
	exception *http.Transport
	hosts     map[string]bool
}

func (t *hostTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if t.hosts[strings.ToLower(r.URL.Hostname())] {
		return t.exception.RoundTrip(r)
	}
	return t.verified.RoundTrip(r)
}

func mirrorClient(jar *cookiejar.Jar, opts CloneOptions) (*http.Client, error) {
	verified := http.DefaultTransport.(*http.Transport).Clone()
	verified.TLSClientConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	if opts.Proxy != "" {
		proxy, err := url.Parse(opts.Proxy)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy: %w", err)
		}
		verified.Proxy = http.ProxyURL(proxy)
	}
	exception := verified.Clone()
	exception.TLSClientConfig = &tls.Config{MinVersion: tls.VersionTLS12, InsecureSkipVerify: true} // Explicit host-scoped opt-in.
	hosts := map[string]bool{}
	for _, host := range opts.InsecureHosts {
		host = strings.ToLower(strings.TrimSpace(host))
		if host == "" || strings.ContainsAny(host, "/:* ") {
			return nil, fmt.Errorf("insecure-host must be an exact hostname: %q", host)
		}
		hosts[host] = true
	}
	client := &http.Client{Transport: &hostTransport{verified, exception, hosts}, Timeout: 30 * time.Second}
	if jar != nil {
		client.Jar = jar
	}
	return client, nil
}
