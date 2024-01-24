package main

import (
	"embed"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	listen = flag.String("listen", ":8080", "Interface:port to listen on")
	index  = flag.String("index", "README.md", "Index (or default) markdown file name")

	//go:embed static
	static embed.FS
)

func render(cwd string, w http.ResponseWriter, r *http.Request) {
	b, err := os.ReadFile(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	css, body, err := renderMarkdown(b)
	fmt.Fprintf(w, `<!doctype html>
	<html lang="en"><head><meta charset="utf-8">
		<title>%s</title>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
		<link href="/static/bootstrap.min.css" rel="stylesheet" type="text/css" />
		<link href="/static/style.css" rel="stylesheet" type="text/css" />
		<style>%s</style>
		<link rel="icon" type="image/png" href="/static/markdown.png" />
	</head><body><article class="markdown-body">`,
		filepath.Base(cwd), css)
	w.Write(body)
	io.WriteString(w, `</article></body></html>`)
}

func markdownRender(cwd string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "" {
			r.URL.Path = *index
		}
		if strings.HasSuffix(r.URL.Path, ".md") {
			render(cwd, w, r)
			log.WithFields(log.Fields{"file": r.URL.Path}).Info("Rendered")
			return
		}
		next.ServeHTTP(w, r)
		log.WithFields(log.Fields{"file": r.URL.Path}).Info("Served")
	})
}

func main() {
	flag.Parse()
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Unable to determine current working directory, err=%s", err)
	}
	http.Handle("/", http.StripPrefix("/", markdownRender(cwd, http.FileServer(http.Dir(cwd)))))
	http.Handle("/static/", http.FileServer(http.FS(static)))

	log.WithFields(log.Fields{"dir": cwd, "addr": *listen}).Info("Starting")
	if err := http.ListenAndServe(*listen, nil); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Failed at runtime")
	}
}
