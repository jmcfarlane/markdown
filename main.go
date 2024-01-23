package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"github.com/shurcooL/github_flavored_markdown"
	"github.com/shurcooL/github_flavored_markdown/gfmstyle"

	log "github.com/sirupsen/logrus"
)

var (
	listen = flag.String("listen", ":8080", "Interface:port to listen on")
	index  = flag.String("index", "README.md", "Index (or default) markdown file name")
)

func render(cwd string, w http.ResponseWriter, r *http.Request) {
	b, err := os.ReadFile(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, `<html><head><meta charset="utf-8">
		<title>%s</title>
		<link href="/assets/gfm.css" media="all" rel="stylesheet" type="text/css" />
		<link href="https://cdnjs.cloudflare.com/ajax/libs/octicons/2.1.2/octicons.css" media="all" rel="stylesheet" type="text/css" />
	</head><body><article class="markdown-body entry-content" style="padding: 30px;">`, filepath.Base(cwd))
	w.Write(github_flavored_markdown.Markdown(b))
	// w.Write(markdownWraper(b))
	io.WriteString(w, `</article></body></html>`)
}

func markdown(cwd string, next http.Handler) http.Handler {
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

// TESTING
func markdownWraper(text []byte) []byte {
	var htmlFlags = 0
	// htmlFlags |= blackfriday.HTML_SKIP_STYLE
	renderer := &renderer{Html: blackfriday.HtmlRenderer(htmlFlags, "", "").(*blackfriday.Html)}
	unsanitized := blackfriday.Markdown(text, renderer, extensions)
	// sanitized := policy.SanitizeBytes(unsanitized)
	return unsanitized
}

type renderer struct {
	*blackfriday.Html
}

// policy for GitHub Flavored Markdown-like sanitization.
var policy = func() *bluemonday.Policy {
	p := bluemonday.UGCPolicy()
	p.AllowAttrs("class").Matching(bluemonday.SpaceSeparatedTokens).OnElements("div", "span")
	p.AllowAttrs("class", "name").Matching(bluemonday.SpaceSeparatedTokens).OnElements("a")
	p.AllowAttrs("rel").Matching(regexp.MustCompile(`^nofollow$`)).OnElements("a")
	p.AllowAttrs("aria-hidden").Matching(regexp.MustCompile(`^true$`)).OnElements("a")
	p.AllowAttrs("type").Matching(regexp.MustCompile(`^checkbox$`)).OnElements("input")
	p.AllowAttrs("checked", "disabled").Matching(regexp.MustCompile(`^$`)).OnElements("input")
	p.AllowDataURIImages()
	return p
}()

// extensions for GitHub Flavored Markdown-like parsing.
const extensions = blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
	blackfriday.EXTENSION_TABLES |
	blackfriday.EXTENSION_FENCED_CODE |
	blackfriday.EXTENSION_AUTOLINK |
	blackfriday.EXTENSION_STRIKETHROUGH |
	blackfriday.EXTENSION_SPACE_HEADERS |
	blackfriday.EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK

// TESTING

func main() {
	flag.Parse()
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Unable to determine current working directory, err=%s", err)
	}
	http.Handle("/", http.StripPrefix("/", markdown(cwd, http.FileServer(http.Dir(cwd)))))
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(gfmstyle.Assets)))
	log.WithFields(log.Fields{"dir": cwd, "addr": *listen}).Info("Listening")
	http.ListenAndServe(*listen, nil)
}
