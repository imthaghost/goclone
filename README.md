<p align="center">
  <a href="https://goclone.io/">
    <img alt="jedi" src="docs/media/logo.png"> 
  </a>
</p>
<p align="center">
Copy websites to your computer! Goclone is a utility that allows you to download a website from the Internet to a local directory. Get html, css, js, images, and other files from the server to your computer. Goclone arranges the original site's relative link-structure. Simply open a page of the "mirrored" website in your browser, and you can browse the site from link to link as if you were viewing it online.
</p>
<br>
<p align="center"><a href="https://goclone.io/">Official Website</a></p>
<br>
<p align="center">
   <a href="https://goreportcard.com/report/github.com/imthaghost/goclone"><img src="https://goreportcard.com/badge/github.com/imthaghost/goclone"></a>
   <a href="https://github.com/imthaghost/goclone/actions/workflows/master-workflow.yml"><img src="https://github.com/imthaghost/goclone/actions/workflows/master-workflow.yml/badge.svg"></a>
   <a href="https://github.com/imthaghost/goclone/blob/master/LICENSE"><img src="https://img.shields.io/badge/License-MIT-yellow.svg"></a>
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
brew tap imthaghost/goclone
# install tool
brew install goclone
```

<a name="manual"></a>

### Manual

```bash
# Go version >= 1.20
go install github.com/imthaghost/goclone/cmd/goclone@latest
```
#### Or

```bash
# go get :)
go get github.com/imthaghost/goclone
# change to project directory using your GOPATH
cd $GOPATH/src/github.com/imthaghost/goclone/cmd/goclone
# build and install application
go install
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

<a name="contributors"></a>

## Contributors

Contributions are welcome! Please see [Contributing Guide](https://github.com/imthaghost/goclone/blob/master/docs/CONTRIBUTING.md) for more details.

<table>
  <tr>
    <td align="center"><a href="https://github.com/imthaghost"><img src="https://avatars3.githubusercontent.com/u/46610773?s=460&v=4" width="75px;" alt="Gary Frederick"/><br /><sub><b>Tha Ghost</b></sub></a><br /><a href="https://github.com/imthaghost/goclone/commits?author=imthaghost" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/imthaghost"><img src="https://avatars.githubusercontent.com/u/29051129?v=4" width="75px;" alt="Juan Mesaglio"/><br /><sub><b>Juan Mesaglio</b></sub></a><br /><a href="https://github.com/mesaglio" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/tempor1s"><img src="https://avatars0.githubusercontent.com/u/29741401?s=460&u=1ca03db5bbb7046bab14f72b7d6e801b9b0ac6f0&v=4" width="75px;" alt="Ben Lafferty"/><br /><sub><b>Ben Lafferty</b></sub></a><br /><a href="https://github.com/imthaghost/goclone/commits?author=tempor1s" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/omarsagoo"><img src="https://avatars3.githubusercontent.com/u/47726951?s=460&u=b806148e1598b97c454820c9c17452db39441177&v=4" width="75px;" alt="Omar Sagoo"/><br /><sub><b>Omar Sagoo</b></sub></a><br /><a href="https://github.com/imthaghost/goclone/commits?author=omarsagoo" title="Code">💻</a></td>
  </tr>
