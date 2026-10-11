package cmd

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMirrorTLSExceptionIsHostScoped(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }))
	defer srv.Close()
	// Clone initializes the standard transport's lazy HTTP/2 configuration.
	_ = http.DefaultTransport.(*http.Transport).Clone()
	defaultTLS := http.DefaultTransport.(*http.Transport).TLSClientConfig
	wasInsecure := defaultTLS != nil && defaultTLS.InsecureSkipVerify
	for _, tc := range []struct {
		hosts   []string
		success bool
	}{
		{nil, false}, {[]string{"example.com"}, false}, {[]string{"127.0.0.1"}, true},
	} {
		client, err := mirrorClient(nil, CloneOptions{InsecureHosts: tc.hosts})
		if err != nil {
			t.Fatal(err)
		}
		resp, err := client.Get(srv.URL)
		if (err == nil) != tc.success {
			t.Fatalf("hosts %v: success=%t err=%v", tc.hosts, tc.success, err)
		}
		if resp != nil {
			resp.Body.Close()
		}
	}
	defaultTLS = http.DefaultTransport.(*http.Transport).TLSClientConfig
	if (defaultTLS != nil && defaultTLS.InsecureSkipVerify) != wasInsecure {
		t.Fatal("changed default transport certificate verification")
	}
	if _, err := mirrorClient(nil, CloneOptions{InsecureHosts: []string{"*"}}); err == nil {
		t.Fatal("accepted wildcard")
	}
}
