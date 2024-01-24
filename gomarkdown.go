package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	mdhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

const TOC = "[TOC]"

var (
	highlightStyle *chroma.Style
	htmlFormatter  *html.Formatter
	renderer       *mdhtml.Renderer
	tocRenderer    *mdhtml.Renderer

	extensions = parser.CommonExtensions |
		parser.Attributes |
		parser.AutoHeadingIDs |
		parser.Mmark |
		parser.NoEmptyLineBeforeBlock |
		parser.OrderedListStart |
		parser.SuperSubscript
)

func init() {
	highlightStyle = mustMakeChromaStyle()
	htmlFormatter = html.New(html.WithClasses(true), html.TabWidth(2))
	renderer = makeRenderer()
	tocRenderer = mdhtml.NewRenderer(
		mdhtml.RendererOptions{Flags: mdhtml.TOC})
}

func makeRenderer() *mdhtml.Renderer {
	htmlFlags := mdhtml.RendererOptions{
		Flags: mdhtml.CommonFlags |
			mdhtml.Smartypants |
			mdhtml.SmartypantsDashes |
			mdhtml.SmartypantsFractions |
			mdhtml.SmartypantsLatexDashes,
		RenderNodeHook: myRenderHook,
	}
	return mdhtml.NewRenderer(htmlFlags)
}

// https://github.com/alecthomas/chroma#supported-languages
func mustMakeChromaStyle() *chroma.Style {
	styleName := "monokailight"
	highlightStyle = styles.Get(styleName)
	if highlightStyle == nil {
		panic(fmt.Sprintf("didn't find style '%s'", styleName))
	}
	return highlightStyle
}

func chromaCSS() ([]byte, error) {
	w := bytes.NewBufferString("")
	htmlFormatter.WriteCSS(w, highlightStyle)
	return w.Bytes(), nil
}

// https://github.com/alecthomas/chroma/blob/master/quick/quick.go
func htmlHighlight(w io.Writer, source, lang, defaultLang string) error {
	if lang == "" {
		lang = defaultLang
	}
	l := lexers.Get(lang)
	if l == nil {
		l = lexers.Analyse(source)
	}
	if l == nil {
		l = lexers.Fallback
	}
	l = chroma.Coalesce(l)

	it, err := l.Tokenise(nil, source)
	if err != nil {
		return err
	}
	return htmlFormatter.Format(w, highlightStyle, it)
}

// an actual rendering of Paragraph is more complicated
func renderCode(w io.Writer, codeBlock *ast.CodeBlock, entering bool) {
	defaultLang := ""
	lang := string(codeBlock.Info)
	htmlHighlight(w, string(codeBlock.Literal), lang, defaultLang)
}

func myRenderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	if code, ok := node.(*ast.CodeBlock); ok {
		renderCode(w, code, entering)
		return ast.GoToNext, true
	}
	return ast.GoToNext, false
}

func renderMarkdown(m []byte) ([]byte, []byte, error) {
	// Parse document (Note: parser is not thread safe, so create each time)
	doc := parser.NewWithExtensions(extensions).Parse(m)

	// Render document
	page := bytes.NewBufferString("")
	page.Write(markdown.Render(doc, renderer))

	// Generate TOC
	header := bytes.NewBufferString("")
	tocRenderer.RenderHeader(header, doc)

	// Write the rendered doc with the TOC replaced
	b := bytes.Replace(page.Bytes(), []byte(TOC), header.Bytes(), 1)

	css, err := chromaCSS()
	return css, b, err
}
