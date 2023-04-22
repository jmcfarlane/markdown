# Markdown

[![Go Report Card](https://goreportcard.com/badge/jmcfarlane/markdown)](https://goreportcard.com/report/jmcfarlane/markdown)

This is a tool to render
[markdown](https://en.wikipedia.org/wiki/Markdown) files locally, for
use with iteration prior to pushing changes to version control.

1. All markdown rendering occurs locally
1. Your browser will fetch some static assets from the internet:
   1. https://cdnjs.cloudflare.com/ajax/libs/octicons/2.1.2/octicons.css
   1. https://cdnjs.cloudflare.com/ajax/libs/octicons/2.1.2/octicons.woff

## Installation

```
$ go install github.com/jmcfarlane/markdown@latest
```

## Usage

```
$ cd ~/where/your/docs/are
$ markdown
```

Point your browser at: http://localhost:8080. Click around, refresh
the page as you make changes.

## Help

```
$ markdown -h
Usage of markdown:
  -index string
    	Index (or default) markdown file name (default "README.md")
  -listen string
    	Interface:port to listen on (default ":8080")
```

## Credit

This program doesn't really do much of anything. All the actual work
is done by:

- https://github.com/russross/blackfriday
- https://github.com/shurcooL/github_flavored_markdown/
- https://golang.org/

As well as things ^ depend upon.

> Inspired by: https://github.com/joeyespo/grip
