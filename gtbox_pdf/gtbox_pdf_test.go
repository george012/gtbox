package gtbox_pdf

import (
	"regexp"
	"strings"
	"testing"
)

func TestCleanPDFText(t *testing.T) {
	got := CleanPDFText("a​ \r\nb\f c\t \n")
	want := "a\nb\n c\n"
	if got != want {
		t.Fatalf("CleanPDFText got %q want %q", got, want)
	}
}

func TestConvertTextToMarkdownGeneric(t *testing.T) {
	raw := strings.Join([]string{
		"Example API Reference",
		"",
		"Service Index",
		"alpha   compute  example.com/alp",
		"ha",
		"General Protocol",
		"",
		"Feature One",
		"",
		"Endpoint 1: POST /v1/sample",
		"Request fields",
		"```",
		"  1   curl -X POST https://example.com/v1/sample \\",
		"  2       -F data=@file.bin",
		"",
		"FAQ",
		"Q: hello",
		"A: world",
	}, "\n") + "\n"

	opts := MarkdownOptions{
		TopHeadings: []string{
			"Example API Reference", "General Protocol", "Feature One", "FAQ",
		},
		SubHeadings:           []string{"Request fields"},
		SectionHeadingPattern: regexp.MustCompile(`^Endpoint[[:space:]]+[0-9]+:[[:space:]]+POST[[:space:]]+/v1/`),
		CodeBlockMarkers:      []string{"```"},
		LineJoinRules: []MarkdownLineJoinRule{
			{Contains: "example.com/alp", NextEquals: "ha", Replace: "example.com/alp", With: "example.com/alpha"},
		},
		Injections: []MarkdownInjection{
			{
				Trigger: "Service Index",
				Heading: "Service Index",
				Lines: []string{
					"| Service | Tier | URL |",
					"|---|---|---|",
					"| alpha | compute | https://example.com/alpha |",
				},
				SkipUntil: "General Protocol",
			},
		},
		QAEnabled: true,
	}

	md := ConvertTextToMarkdown(raw, opts)
	for _, want := range []string{
		"# Example API Reference",
		"# Service Index",
		"| alpha | compute | https://example.com/alpha |",
		"# General Protocol",
		"# Feature One",
		"## Endpoint 1: POST /v1/sample",
		"### Request fields",
		"curl -X POST https://example.com/v1/sample \\",
		"**Q:** hello",
		"**A:** world",
	} {
		if !strings.Contains(md, want) {
			t.Fatalf("markdown missing %q\n----\n%s", want, md)
		}
	}

	if strings.Contains(md, "example.com/alp\nha") {
		t.Fatalf("LineJoinRules failed to repair the broken URL:\n%s", md)
	}

	report := ReportMarkdown(md, MarkdownReportOptions{
		SectionHeadingPattern: regexp.MustCompile(`^## Endpoint[[:space:]]+[0-9]+:`),
	})
	if report.FormFeedCount != 0 {
		t.Fatalf("form feed count: got %d", report.FormFeedCount)
	}
	if report.CodeFenceCount%2 != 0 {
		t.Fatalf("unbalanced fences: %d", report.CodeFenceCount)
	}
	if report.SectionHeadingCount != 1 {
		t.Fatalf("section headings: got %d", report.SectionHeadingCount)
	}
}

func TestConvertTextToMarkdownNumberedCodeBlockAutoDetect(t *testing.T) {
	raw := strings.Join([]string{
		"Intro line",
		"  1   import os",
		"  2   print(\"hi\")",
		"After block",
	}, "\n") + "\n"

	md := ConvertTextToMarkdown(raw, MarkdownOptions{})
	if !strings.Contains(md, "```\nimport os") {
		t.Fatalf("expected auto-detected fenced block:\n%s", md)
	}
	if strings.Count(md, "```")%2 != 0 {
		t.Fatalf("unbalanced fences:\n%s", md)
	}
}

func TestReportMarkdownDefaults(t *testing.T) {
	md := strings.Join([]string{
		"# Title",
		"",
		"## Section A",
		"some prose",
		"## Section B",
		"```",
		"  1   code line",
		"```",
		"trailing",
	}, "\n") + "\n"

	report := ReportMarkdown(md, MarkdownReportOptions{})
	if report.SectionHeadingCount != 2 {
		t.Fatalf("section headings: got %d", report.SectionHeadingCount)
	}
	if report.CodeFenceCount%2 != 0 {
		t.Fatalf("unbalanced fences: %d", report.CodeFenceCount)
	}
	if report.FormFeedCount != 0 {
		t.Fatalf("form feeds: got %d", report.FormFeedCount)
	}
	if len(report.NumberedLinesOutside) != 0 {
		t.Fatalf("numbered lines outside code should be empty, got %v", report.NumberedLinesOutside)
	}
}
