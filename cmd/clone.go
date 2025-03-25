package cmd

import (
	"context"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"runtime"
	"strings"

	"os/exec"

	"github.com/imthaghost/goclone/pkg/crawler"
	"github.com/imthaghost/goclone/pkg/file"
	"github.com/imthaghost/goclone/pkg/parser"
	"github.com/imthaghost/goclone/pkg/server"
)

// Clone the given site :)
func cloneSite(ctx context.Context, args []string) error {
	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		return err
	}
	var cs []*http.Cookie
	if len(Cookies) != 0 {
		cs = make([]*http.Cookie, 0, len(Cookies))
		for _, c := range Cookies {
			ff := strings.Fields(c)
			for _, f := range ff {
				var k, v string
				if i := strings.IndexByte(f, '='); i >= 0 {
					k, v = f[:i], strings.TrimRight(f[i+1:], ";")
				} else {
					return fmt.Errorf("No = in cookie %q", c)
				}
				cs = append(cs, &http.Cookie{Name: k, Value: v})
			}
		}
		for _, a := range args {
			u, err := url.Parse(a)
			if err != nil {
				return fmt.Errorf("%q: %w", a, err)
			}
			jar.SetCookies(&url.URL{Scheme: u.Scheme, User: u.User, Host: u.Host}, cs)
		}
	}

	var firstProject string
	for _, u := range args {
		isValid, isValidDomain := parser.ValidateURL(u), parser.ValidateDomain(u)
		if !isValid && !isValidDomain {
			return fmt.Errorf("%q is not valid", u)
		}
		name := u
		if isValidDomain {
			u = parser.CreateURL(name)
		} else {
			name = parser.GetDomain(u)
		}
		projectPath := file.CreateProject(name)
		if firstProject == "" {
			firstProject = projectPath
		}

		if err := crawler.Crawl(ctx, u, projectPath, jar, ProxyString, UserAgent, Referer, Depth); err != nil {
			return fmt.Errorf("%q: %w", u, err)
		}
		// Restructure html
		//if err := html.LinkRestructure(projectPath); err != nil {
		//	return fmt.Errorf("%q: %w", projectPath, err)
		//}

	}
	if Serve {
		serverUrl := fmt.Sprintf("http://localhost:%d", ServePort)
		cmd := exec.Command("open", serverUrl)
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("%v: %w", cmd.Args, err)
		}
		return server.Serve(firstProject, ServePort)
	} else if Open {
		// automatically open project
		cmd := open(firstProject + "/index.html")
		if err := cmd.Start(); err != nil {
			return fmt.Errorf("%v: %w", cmd.Args, err)
		}
	}
	return nil
}

// open opens the specified URL in the default browser of the user.
func open(url string) *exec.Cmd {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...)
}
