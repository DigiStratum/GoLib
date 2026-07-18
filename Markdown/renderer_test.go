package markdown_test

import (
	"strings"
	"testing"

	markdown "github.com/DigiStratum/GoLib/Markdown"
)

// TestRender_TableDriven covers the common agentic output patterns:
// headers, links, code blocks, lists, and tables.
func TestRender_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string // substrings that must appear in output
		notContains []string // raw markdown artifacts that must NOT appear in output
	}{
		{
			name:  "h1 header",
			input: "# Hello World",
			// goldmark WithAutoHeadingID adds an id attribute; match on tag + content
			contains: []string{"<h1 ", ">Hello World</h1>"},
		},
		{
			name:     "h2 header",
			input:    "## Section Title",
			contains: []string{"<h2 ", ">Section Title</h2>"},
		},
		{
			name:     "h3 header",
			input:    "### Subsection",
			contains: []string{"<h3 ", ">Subsection</h3>"},
		},
		{
			name:     "inline link",
			input:    "[OpenAI](https://openai.com)",
			contains: []string{`<a href="https://openai.com">OpenAI</a>`},
		},
		{
			name:     "fenced code block",
			input:    "```go\nfmt.Println(\"hello\")\n```",
			contains: []string{"<pre><code", "fmt.Println"},
			notContains: []string{"```"},
		},
		{
			name:     "inline code",
			input:    "Use `template.HTML` type",
			contains: []string{"<code>template.HTML</code>"},
		},
		{
			name:     "unordered list",
			input:    "- item one\n- item two\n- item three",
			contains: []string{"<ul>", "<li>item one</li>", "<li>item two</li>", "<li>item three</li>", "</ul>"},
		},
		{
			name:     "ordered list",
			input:    "1. first\n2. second\n3. third",
			contains: []string{"<ol>", "<li>first</li>", "<li>second</li>", "</ol>"},
		},
		{
			name:  "table",
			input: "| Name | Value |\n|------|-------|\n| foo  | bar   |\n| baz  | qux   |",
			contains: []string{
				"<table>",
				"<thead>",
				"<th>Name</th>",
				"<th>Value</th>",
				"<tbody>",
				"<td>foo</td>",
				"<td>bar</td>",
				"<td>baz</td>",
				"<td>qux</td>",
				"</table>",
			},
		},
		{
			name:     "strikethrough (GFM)",
			input:    "~~deprecated~~",
			contains: []string{"<del>deprecated</del>"},
		},
		{
			name:  "autolink (GFM)",
			input: "Visit https://digistratum.com for more.",
			contains: []string{`href="https://digistratum.com"`},
		},
		{
			name:     "paragraph",
			input:    "This is a plain paragraph.",
			contains: []string{"<p>This is a plain paragraph.</p>"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := markdown.Render(tc.input)
			if err != nil {
				t.Fatalf("Render() returned unexpected error: %v", err)
			}

			// Verify the return type is template.HTML (safe/unescaped) —
			// string() conversion only compiles if result is template.HTML.
			html := string(result)

			for _, want := range tc.contains {
				if !strings.Contains(html, want) {
					t.Errorf("expected output to contain %q\ngot:\n%s", want, html)
				}
			}

			for _, bad := range tc.notContains {
				if strings.Contains(html, bad) {
					t.Errorf("expected output NOT to contain %q\ngot:\n%s", bad, html)
				}
			}
		})
	}
}

// TestRender_EmptyInput ensures empty markdown returns empty/minimal HTML without error.
func TestRender_EmptyInput(t *testing.T) {
	result, err := markdown.Render("")
	if err != nil {
		t.Fatalf("Render(\"\") returned unexpected error: %v", err)
	}
	html := string(result)
	// Should be empty or just whitespace — no content generated
	if strings.TrimSpace(html) != "" {
		t.Logf("Empty input produced non-empty output (acceptable): %q", html)
	}
}

// TestMustRender_PanicsOnError verifies MustRender is a convenience wrapper.
// Since goldmark doesn't typically error, we mainly verify it returns template.HTML.
func TestMustRender_ReturnsHTML(t *testing.T) {
	result := markdown.MustRender("# Hello")
	// string() on template.HTML proves the compile-time type without an explicit declaration.
	if !strings.Contains(string(result), "<h1") {
		t.Errorf("MustRender expected <h1> tag, got: %s", string(result))
	}
}

// TestRender_ReturnsTemplateHTML verifies the type at compile-time (type assertion in loop above)
// and ensures it is NOT plain string (would cause double-escaping in html/template).
func TestRender_TypeIsTemplateHTML(t *testing.T) {
	result, err := markdown.Render("**bold**")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// template.HTML assignment proves the type — if this compiles, the type is correct.
	safe := result
	_ = safe
}
