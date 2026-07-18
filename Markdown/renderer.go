// Package markdown provides a thin goldmark-based Markdown → HTML renderer for
// use across DigiStratum Go applications.
//
// The primary entry-point is [Render], which converts a Markdown string to a
// [html/template.HTML] value.  Because [html/template.HTML] is treated as
// pre-sanitised by Go's html/template engine, the rendered output is inserted
// verbatim into templates — no double-escaping occurs.
//
// GFM extensions enabled by default:
//   - Tables
//   - Strikethrough
//   - Autolinks
//   - Task lists
//
// Safe HTML mode is ON: raw HTML blocks inside the Markdown source are
// stripped, protecting downstream templates from injection.
package markdown

import (
	"bytes"
	"html/template"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	gmHTML "github.com/yuin/goldmark/renderer/html"
)

// defaultParser is a package-level goldmark instance with sensible defaults
// for DigiStratum agentic output patterns.  It is safe for concurrent use.
var defaultParser = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM, // Tables, Strikethrough, Autolinks, TaskList
	),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(), // stable heading anchors for in-page links
	),
	goldmark.WithRendererOptions(
		gmHTML.WithHardWraps(), // treat single newlines as <br>
		gmHTML.WithXHTML(),     // well-formed XHTML output (self-closing tags)
	),
)

// Render converts md (a Markdown string) into HTML and returns it as a
// [html/template.HTML] value, which is safe to use directly in Go templates
// without further escaping.
//
// The renderer uses GitHub Flavoured Markdown extensions (tables, strikethrough,
// autolinks, task lists) and strips any raw HTML present in the source so that
// agentic output cannot inject markup.
//
// Example:
//
//	html, err := markdown.Render("## Hello\n\n- item one\n- item two")
//	// html == template.HTML("<h2 id=\"hello\">Hello</h2>\n<ul>\n<li>item one</li>\n…")
func Render(md string) (template.HTML, error) {
	var buf bytes.Buffer
	if err := defaultParser.Convert([]byte(md), &buf); err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil //nolint:gosec // caller-controlled markdown; raw HTML stripped by goldmark safe mode
}

// MustRender is like [Render] but panics on error.  Use in init() or tests
// where the input is a compile-time constant.
func MustRender(md string) template.HTML {
	html, err := Render(md)
	if err != nil {
		panic("markdown.MustRender: " + err.Error())
	}
	return html
}
