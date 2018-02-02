package main

import (
	"flag"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/shurcooL/github_flavored_markdown"
	"github.com/shurcooL/github_flavored_markdown/gfmstyle"

	log "github.com/sirupsen/logrus"
)

var (
	listen = flag.String("listen", ":8080", "Interface:port to listen on")
	index  = flag.String("index", "README.md", "Index (or default) markdown file name")
)

func render(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadFile(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	io.WriteString(w, `<html><head><meta charset="utf-8">
		<link href="/assets/gfm.css" media="all" rel="stylesheet" type="text/css" />
		<link href="https://cdnjs.cloudflare.com/ajax/libs/octicons/2.1.2/octicons.css" media="all" rel="stylesheet" type="text/css" />
	</head><body><article class="markdown-body entry-content" style="padding: 30px;">`)
	w.Write(github_flavored_markdown.Markdown(b))
	io.WriteString(w, `</article></body></html>`)
}

func markdown(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "" {
			r.URL.Path = *index
		}
		if strings.HasSuffix(r.URL.Path, ".md") {
			render(w, r)
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
	http.Handle("/", http.StripPrefix("/", markdown(http.FileServer(http.Dir(cwd)))))
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(gfmstyle.Assets)))
	log.WithFields(log.Fields{"dir": cwd, "addr": *listen}).Info("Listening")
	http.ListenAndServe(*listen, nil)
}
