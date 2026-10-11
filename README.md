<p align="center">
  <a href="https://goclone.io/">
    <img alt="jedi" src="docs/media/logo.png"> 
  </a>
</p>
<p align="center">
Copy websites to your computer! Goclone is a utility that allows you to download a website from the Internet to a local directory. Get html, css, js, images, and other files from the server to your computer. Goclone arranges the original site's relative link-structure. Simply open a page of the "mirrored" website in your browser, and you can browse the site from link to link as if you were viewing it online.
</p>
<br>
<p align="center"><a href="https://goclone.dev/">Official Website</a></p>
<br>
<p align="center">
   <a href="https://goreportcard.com/report/github.com/goclone-dev/goclone"><img src="https://goreportcard.com/badge/github.com/goclone-dev/goclone"></a>
   <a href="https://github.com/goclone-dev/goclone/actions/workflows/master-workflow.yml"><img src="https://github.com/goclone-dev/goclone/actions/workflows/master-workflow.yml/badge.svg"></a>
   <a href="https://github.com/goclone-dev/goclone/blob/master/LICENSE"><img src="https://img.shields.io/badge/License-MIT-yellow.svg"></a>
</p>
<br>

![Example](/docs/media/bitski.gif)

## Table of Contents

- [Installation](#installation)
  - [Brew](#brew)
  - [Manual](#manual)
- [Examples](#examples)
- [Contributors](#contributors)

<a name="installation"></a>

## 🚀 Installation

<a name="brew"></a>

### Brew

```bash
# tap
brew tap goclone-dev/goclone
# install tool
brew install goclone
```

<a name="manual"></a>

### Manual

```bash
# Go version >= 1.20
go install github.com/goclone-dev/goclone/cmd/goclone@latest
```

#### Or

```bash
# Clone the repository
git clone https://github.com/goclone-dev/goclone.git
cd goclone
# Build and run
go build -o goclone cmd/goclone/main.go
# Move binary to a directory in your PATH (optional)
mv goclone /usr/local/bin/
```

<a name="examples"></a>

## Examples

```bash
# goclone <url>
goclone https://configtree.co
```

![Config](/docs/media/config.gif)

## Usage

```
Usage:
  goclone <url> [flags]

Flags:
  -C, --cookie strings        Pre-set these cookies
  -h, --help                  help for goclone
  -o, --open                  Automatically open project in default browser
  -p, --proxy_string string   Proxy connection string. Support http and socks5 https://pkg.go.dev/github.com/gocolly/colly#Collector.SetProxy
  -s, --serve                 Serve the generated files using Echo.
  -P, --servePort int         Serve port number. (default 5000)
  -u, --user_agent string     Custom User Agent
```

## Mirroring WordPress and other static sites

Use the opt-in mirror mode when the site depends on responsive images, fonts,
CSS imports/backgrounds, or several linked pages:

```sh
goclone --mirror --delay 500ms --max-pages 50 https://example.com/
```

Mirror mode preserves directory paths, downloads HTML and CSS dependencies
(including WebP/AVIF, `srcset`, common lazy-load attributes, CSS `url()` and
`@import`), and follows same-origin page links. External assets are placed in
separate origin directories. Query strings remain part of the HTTP request and
are hashed into local filenames to prevent collisions and invalid `?` filenames.
References are rewritten to local paths. Serve the output directory over HTTP;
root-relative URLs are not intended for `file://` browsing.

Use repeatable `--page /path/` flags to include pages found in a sitemap but not
linked from the entry page. The total includes these seeds and respects
`--max-pages`. The crawler does not submit forms or follow admin/API routes,
query-string navigation links, or off-origin page links. Every request uses the
configured user agent, cookies, proxy, timeout, and cancellation context. Requests
are sequential with `--delay` between them. A resource is limited to 64 MiB.

`mirror-manifest.json` records source URLs, local files, MIME types, and failures.
A failed download produces a nonzero exit status rather than silent success.
The legacy single-page mode remains unchanged.

TLS certificates are verified by default in mirror mode. For a known expired
certificate on a site you administer, `--insecure-host example.com` provides an
explicit exception for that exact hostname; redirects to other hosts retain
verification. This option does not modify the global HTTP transport.

This captures the public frontend, not WordPress's PHP source or database.
JavaScript is not executed, so runtime-generated URLs, webpack chunks, embedded
applications, and server-backed forms may need additional capture or a backend.
The CSS scanner supports ordinary URL/import syntax, not a complete CSS parser;
mixed data-URL `srcset` values are preserved without processing their other entries.

<a name="contributors"></a>

## Contributors

Contributions are welcome! Please see [Contributing Guide](https://github.com/goclone-dev/goclone/blob/master/docs/CONTRIBUTING.md) for more details.

<table>
  <tr>
    <td align="center"><a href="https://github.com/imthaghost"><img src="https://avatars3.githubusercontent.com/u/46610773?s=460&v=4" width="75px;" alt="Gary Frederick"/><br /><sub><b>Tha Ghost</b></sub></a><br /><a href="https://github.com/goclone-dev/goclone/commits?author=imthaghost" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/imthaghost"><img src="https://avatars.githubusercontent.com/u/29051129?v=4" width="75px;" alt="Juan Mesaglio"/><br /><sub><b>Juan Mesaglio</b></sub></a><br /><a href="https://github.com/mesaglio" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/tempor1s"><img src="https://avatars0.githubusercontent.com/u/29741401?s=460&u=1ca03db5bbb7046bab14f72b7d6e801b9b0ac6f0&v=4" width="75px;" alt="Ben Lafferty"/><br /><sub><b>Ben Lafferty</b></sub></a><br /><a href="https://github.com/goclone-dev/goclone/commits?author=tempor1s" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/omarsagoo"><img src="https://avatars3.githubusercontent.com/u/47726951?s=460&u=b806148e1598b97c454820c9c17452db39441177&v=4" width="75px;" alt="Omar Sagoo"/><br /><sub><b>Omar Sagoo</b></sub></a><br /><a href="https://github.com/goclone-dev/goclone/commits?author=omarsagoo" title="Code">💻</a></td>
  </tr>
