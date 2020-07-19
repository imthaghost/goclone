<p align="center">
  <a href="https://goclone.herokuapp.com/">
    <img alt="jedi" src="docs/media/logo.png"> 
  </a>
</p>
<p align="center">
Copy websites to your computer! goclone is a utility that allows you to download a website from the Internet to a local directory. Get html, css, js, images, and other files from the server to your computer. goclone arranges the original site's relative link-structure. Simply open a page of the "mirrored" website in your browser, and you can browse the site from link to link as if you were viewing it online.
</p>
<br>
<p align="center"><a href="https://goclone.app.imthaghost.dev/">Official Website</a></p>
<br>
<p align="center">
   <a href="https://goreportcard.com/report/github.com/imthaghost/goclone"><img src="https://goreportcard.com/badge/github.com/imthaghost/goclone"></a>
   <a href="https://travis-ci.org/imthaghost/goclone.svg?branch=master"><img src="https://travis-ci.org/imthaghost/goclone.svg?branch=master"></a>
   <a href="https://github.com/imthaghost/gitmoji-changelog"><img src="https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg"alt="gitmoji-changelog"></a>
</p>
<br>

![Example](/docs/media/fast.gif)

## Table of Contents

-   [Installation](#installation)
-   [Todo](#Todo)
-   [Examples](#examples)
-   [License](#license)
-   [Contributors](#contributors)

## 🚀 Installation

### Brew

```bash
# tap
brew tap imthaghost/goclone
# install tool
brew install goclone
```

### Manual

```bash
# go get :)
go get github.com/imthaghost/goclone
# change to project directory using your GOPATH
cd $GOPATH/src/github.com/imthaghost/goclone
# build and install application
go install
```

## Todo

### Short term

-   [x] Clone top level site only
-   [x] Update command line interface with Cobra
-   [ ] Clone all pages with given domain
-   [ ] 80-100% test coverage
-   [ ] Update scraper for better performance

### Long term
-   [ ] Clone site that sits behind authentication wall
-   [ ] User specified depth of clone

## Examples

###

```bash
# goclone <url>
goclone https://dribbble.com
```

![Dribbble](/docs/media/dribbble.gif)

## 📝 License

By contributing, you agree that your contributions will be licensed under its MIT License.

In short, when you submit code changes, your submissions are understood to be under the same [MIT License](http://choosealicense.com/licenses/mit/) that covers the project. Feel free to contact the maintainers if that's a concern.

## Contributors

Contributions are welcome! Please see [Contributing Guide](https://github.com/imthaghost/goclone/blob/master/docs/CONTRIBUTING.md) for more details.

<table>
  <tr>
    <td align="center"><a href="https://github.com/imthaghost"><img src="https://avatars3.githubusercontent.com/u/46610773?s=460&v=4" width="75px;" alt="Gary Frederick"/><br /><sub><b>Tha Ghost</b></sub></a><br /><a href="https://github.com/imthaghost/goclone/commits?author=imthaghost" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/tempor1s"><img src="https://avatars0.githubusercontent.com/u/29741401?s=460&u=1ca03db5bbb7046bab14f72b7d6e801b9b0ac6f0&v=4" width="75px;" alt="Ben Lafferty"/><br /><sub><b>Ben Lafferty</b></sub></a><br /><a href="https://github.com/imthaghost/goclone/commits?author=tempor1s" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/omarsagoo"><img src="https://avatars3.githubusercontent.com/u/47726951?s=460&u=b806148e1598b97c454820c9c17452db39441177&v=4" width="75px;" alt="Omar Sagoo"/><br /><sub><b>Omar Sagoo</b></sub></a><br /><a href="https://github.com/imthaghost/goclone/commits?author=omarsagoo" title="Code">💻</a></td>
  </tr>
