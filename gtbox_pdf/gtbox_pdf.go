/*
Package gtbox_pdf is the PDF utility namespace for gtbox.

Current capabilities:
  - text extraction (pure Go default; external pdftotext is an opt-in adapter)
  - text → Markdown structuring (data-driven; the SDK carries no business
    rules)

Planned capabilities (this package is the namespace they will land in):
metadata, split / merge, image extraction, watermark, compression, signature
validation, PDF → HTML.

Design rules:
  - The default code path is pure Go and cross-platform (no shell, awk, sed,
    pdftotext, etc.).
  - The SDK never holds caller-specific data: headings, URL repairs, injected
    blocks and similar must be supplied by the caller. Business documents
    keep their rule sets in their own repositories or load them from
    external configuration files at the CLI boundary.
*/
package gtbox_pdf

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	numberedLineRe = regexp.MustCompile(`^[[:space:]]*[0-9]+[[:space:]]+`)
	loneNumberRe   = regexp.MustCompile(`^[[:space:]]*[0-9]+[[:space:]]*$`)
)

// CleanPDFText normalises line endings, strips zero-width spaces, and turns
// form-feed page separators into newlines. It also right-trims each line.
//
// This is the canonical pre-processing step before ConvertTextToMarkdown; the
// MarkdownReport invariants (form_feeds=0 in particular) depend on it.
func CleanPDFText(raw string) string {
	s := strings.ReplaceAll(raw, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	s = strings.ReplaceAll(s, "​", "")
	s = strings.ReplaceAll(s, "\f", "\n")

	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRightFunc(line, unicode.IsSpace)
	}
	return strings.Join(lines, "\n")
}
