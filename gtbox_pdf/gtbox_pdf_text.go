package gtbox_pdf

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/ledongthuc/pdf"
)

// ExtractText returns the textual content of a PDF file using the SDK's
// default pure-Go engine. The implementation walks pages in document order
// and groups runs by their Y-coordinate row so the output keeps a
// human-readable visual order without depending on any external binary.
//
// Backing library: github.com/ledongthuc/pdf (BSD-3-Clause).
//
// The function honours ctx between pages; the underlying parser is itself
// synchronous, so a very large page can outrun a tight deadline.
func ExtractText(ctx context.Context, pdfPath string) (string, error) {
	if strings.TrimSpace(pdfPath) == "" {
		return "", errors.New("gtbox_pdf: empty pdf path")
	}
	if err := ctx.Err(); err != nil {
		return "", err
	}

	f, r, err := pdf.Open(pdfPath)
	if err != nil {
		return "", fmt.Errorf("gtbox_pdf: open pdf: %w", err)
	}
	defer f.Close()

	var buf bytes.Buffer
	totalPages := r.NumPage()
	for pageIndex := 1; pageIndex <= totalPages; pageIndex++ {
		if err := ctx.Err(); err != nil {
			return "", err
		}
		page := r.Page(pageIndex)
		if page.V.IsNull() {
			continue
		}
		rows, err := page.GetTextByRow()
		if err != nil {
			return "", fmt.Errorf("gtbox_pdf: page %d: %w", pageIndex, err)
		}
		for _, row := range rows {
			var line strings.Builder
			var prev pdf.Text
			for j, w := range row.Content {
				if j > 0 && needsSpaceBetween(prev, w) {
					line.WriteByte(' ')
				}
				line.WriteString(w.S)
				prev = w
			}
			buf.WriteString(line.String())
			buf.WriteByte('\n')
		}
		if pageIndex < totalPages {
			buf.WriteByte('\n')
		}
	}
	return buf.String(), nil
}

// needsSpaceBetween inserts a space when two adjacent text runs are visually
// separated on the page. The library returns a run per character, so without
// this the output collapses into a continuous string. The 0.3 * font-size
// gap is a conservative heuristic; smaller gaps belong to the same word.
func needsSpaceBetween(prev, next pdf.Text) bool {
	if prev.S == "" {
		return false
	}
	gap := next.X - (prev.X + prev.W)
	threshold := 0.3 * next.FontSize
	if threshold <= 0 {
		threshold = 1.0
	}
	return gap > threshold
}

// ExtractTextWithPdftotext is an explicit opt-in adapter that shells out to
// the external `pdftotext` binary (Poppler). It exists for callers who need
// pdftotext's `-layout` output specifically; the SDK never selects this path
// implicitly.
//
// The caller is responsible for ensuring pdftotext is installed and for
// understanding its license (Poppler is GPL-2.0); shipping pdftotext as a
// runtime dependency is the caller's choice.
func ExtractTextWithPdftotext(ctx context.Context, pdfPath string) (string, error) {
	if strings.TrimSpace(pdfPath) == "" {
		return "", errors.New("gtbox_pdf: empty pdf path")
	}
	cmd := exec.CommandContext(ctx, "pdftotext", "-enc", "UTF-8", "-layout", pdfPath, "-")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("gtbox_pdf: pdftotext failed: %s", msg)
	}
	return stdout.String(), nil
}
