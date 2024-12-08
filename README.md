# Markdown

[![Go Report Card](https://goreportcard.com/badge/jmcfarlane/markdown)](https://goreportcard.com/report/jmcfarlane/markdown)

This is a tool to render
[markdown](https://en.wikipedia.org/wiki/Markdown) files locally, for
use with iteration prior to pushing changes to version control.

1. All markdown rendering occurs locally

## Installation

```sh
$ go install github.com/jmcfarlane/markdown@latest
```

## Usage

```sh
$ cd ~/where/your/docs/are
$ markdown
```

Point your browser at: http://localhost:8080. Click around, refresh
the page as you make changes.

## Help

```sh
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

- [github.com/alecthomas/chroma](https://github.com/alecthomas/chroma)
- [github.com/gomarkdown/markdown](https://github.com/gomarkdown/markdown)
- [golang.org](https://golang.org)

As well as things ^ depend upon.

> Inspired by: https://github.com/joeyespo/grip
