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

var extensions = parser.CommonExtensions |
	parser.Attributes |
	parser.AutoHeadingIDs |
	parser.Mmark |
	parser.NoEmptyLineBeforeBlock |
	parser.OrderedListStart |
	parser.SuperSubscript

func makeRenderer(formatter *html.Formatter, style *chroma.Style) *mdhtml.Renderer {
	htmlFlags := mdhtml.RendererOptions{
		Flags: mdhtml.CommonFlags |
			mdhtml.Smartypants |
			mdhtml.SmartypantsDashes |
			mdhtml.SmartypantsFractions |
			mdhtml.SmartypantsLatexDashes,
		RenderNodeHook: myRenderHook(formatter, style),
	}
	return mdhtml.NewRenderer(htmlFlags)
}

// https://github.com/alecthomas/chroma#supported-languages
func makeChromaStyle() (*chroma.Style, error) {
	styleName := "monokai"
	style := styles.Get(styleName)
	if style == nil {
		return nil, fmt.Errorf("Failed to get style=%s", styleName)
	}
	return style, nil
}

// https://github.com/alecthomas/chroma/blob/master/quick/quick.go
func htmlHighlight(w io.Writer, formatter *html.Formatter, style *chroma.Style, source, lang, defaultLang string) error {
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
	return formatter.Format(w, style, it)
}

// an actual rendering of Paragraph is more complicated
func renderCode(w io.Writer, formatter *html.Formatter, style *chroma.Style, codeBlock *ast.CodeBlock, entering bool) {
	defaultLang := ""
	lang := string(codeBlock.Info)
	htmlHighlight(w, formatter, style, string(codeBlock.Literal), lang, defaultLang)
}

func myRenderHook(formatter *html.Formatter, style *chroma.Style) mdhtml.RenderNodeFunc {
	return func(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
		if code, ok := node.(*ast.CodeBlock); ok {
			renderCode(w, formatter, style, code, entering)
			return ast.GoToNext, true
		}
		return ast.GoToNext, false
	}
}

func renderMarkdown(m []byte) ([]byte, []byte, error) {
	// Parse markdown (Note: bits are not thread safe, so create each time)
	myFormatter := html.New(html.WithClasses(true), html.TabWidth(2))
	myParser := parser.NewWithExtensions(extensions).Parse(m)

	// Extract chroma style
	style, err := makeChromaStyle()
	if err != nil {
		return nil, nil, err
	}

	// Render body
	body := bytes.NewBufferString("")
	body.Write(markdown.Render(myParser, makeRenderer(myFormatter, style)))

	// Render CSS
	css := bytes.NewBufferString("")
	myFormatter.WriteCSS(css, style)

	// Render TOC
	toc := bytes.NewBufferString("")
	mdhtml.NewRenderer(mdhtml.RendererOptions{Flags: mdhtml.TOC}).RenderHeader(toc, myParser)

	// Inject the TOC into the body
	b := bytes.Replace(body.Bytes(), []byte(TOC), toc.Bytes(), 1)

	// Render css and body
	return css.Bytes(), b, nil
}
