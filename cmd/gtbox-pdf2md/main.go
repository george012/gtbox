package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/george012/gtbox/gtbox_pdf"
)

// configFile mirrors gtbox_pdf.MarkdownOptions but expresses regex fields as
// strings so the rule set can live in a JSON file outside the SDK. The CLI is
// the only consumer; the SDK itself never reads configuration files.
type configFile struct {
	TopHeadings           []string                         `json:"top_headings"`
	SubHeadings           []string                         `json:"sub_headings"`
	SectionHeadingPattern string                           `json:"section_heading_pattern"`
	CodeBlockMarkers      []string                         `json:"code_block_markers"`
	LineJoinRules         []gtbox_pdf.MarkdownLineJoinRule `json:"line_join_rules"`
	Injections            []gtbox_pdf.MarkdownInjection    `json:"injections"`
	QAEnabled             bool                             `json:"qa_enabled"`
	SuspiciousLinePattern string                           `json:"suspicious_line_pattern"`
}

type extractor func(ctx context.Context, path string) (string, error)

func main() {
	in := flag.String("in", "", "input PDF path")
	out := flag.String("out", "", "output Markdown path")
	configPath := flag.String("config", "", "optional JSON config supplying MarkdownOptions rules (caller-provided)")
	extractorName := flag.String("extractor", "go", "text extractor: go (pure-Go default) | pdftotext (external opt-in)")
	timeout := flag.Duration("timeout", 2*time.Minute, "extractor timeout")
	reportOnly := flag.Bool("report", false, "print report for an existing markdown file")
	flag.Parse()

	if *out == "" {
		fatalf("-out is required")
	}

	mdOpts, reportOpts, err := loadConfig(*configPath)
	if err != nil {
		fatalf("load config: %v", err)
	}

	if *reportOnly {
		data, err := os.ReadFile(*out)
		if err != nil {
			fatalf("read markdown: %v", err)
		}
		fmt.Println(gtbox_pdf.ReportMarkdown(string(data), reportOpts).String())
		return
	}

	if *in == "" {
		fatalf("-in is required")
	}

	extract, err := pickExtractor(*extractorName)
	if err != nil {
		fatalf("%v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	text, err := extract(ctx, *in)
	if err != nil {
		fatalf("%v", err)
	}

	md := gtbox_pdf.ConvertTextToMarkdown(text, mdOpts)
	if err := os.WriteFile(*out, []byte(md), 0644); err != nil {
		fatalf("write markdown: %v", err)
	}

	fmt.Println(gtbox_pdf.ReportMarkdown(md, reportOpts).String())
}

func pickExtractor(name string) (extractor, error) {
	switch strings.TrimSpace(name) {
	case "", "go":
		return gtbox_pdf.ExtractText, nil
	case "pdftotext":
		return gtbox_pdf.ExtractTextWithPdftotext, nil
	default:
		return nil, fmt.Errorf("unknown extractor %q (expected: go | pdftotext)", name)
	}
}

func loadConfig(path string) (gtbox_pdf.MarkdownOptions, gtbox_pdf.MarkdownReportOptions, error) {
	if strings.TrimSpace(path) == "" {
		return gtbox_pdf.MarkdownOptions{}, gtbox_pdf.MarkdownReportOptions{}, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return gtbox_pdf.MarkdownOptions{}, gtbox_pdf.MarkdownReportOptions{}, err
	}
	var cfg configFile
	if err := json.Unmarshal(data, &cfg); err != nil {
		return gtbox_pdf.MarkdownOptions{}, gtbox_pdf.MarkdownReportOptions{}, fmt.Errorf("parse %s: %w", path, err)
	}

	mdOpts := gtbox_pdf.MarkdownOptions{
		TopHeadings:      cfg.TopHeadings,
		SubHeadings:      cfg.SubHeadings,
		CodeBlockMarkers: cfg.CodeBlockMarkers,
		LineJoinRules:    cfg.LineJoinRules,
		Injections:       cfg.Injections,
		QAEnabled:        cfg.QAEnabled,
	}
	if cfg.SectionHeadingPattern != "" {
		re, err := regexp.Compile(cfg.SectionHeadingPattern)
		if err != nil {
			return gtbox_pdf.MarkdownOptions{}, gtbox_pdf.MarkdownReportOptions{}, fmt.Errorf("section_heading_pattern: %w", err)
		}
		mdOpts.SectionHeadingPattern = re
	}

	reportOpts := gtbox_pdf.MarkdownReportOptions{}
	if cfg.SectionHeadingPattern != "" {
		// Reuse the same regex but anchored after the "## " prefix the
		// converter emits, so the report counts the headings the converter
		// just produced.
		re, err := regexp.Compile(`^## ` + trimAnchor(cfg.SectionHeadingPattern))
		if err != nil {
			return gtbox_pdf.MarkdownOptions{}, gtbox_pdf.MarkdownReportOptions{}, fmt.Errorf("report section heading regex: %w", err)
		}
		reportOpts.SectionHeadingPattern = re
	}
	if cfg.SuspiciousLinePattern != "" {
		re, err := regexp.Compile(cfg.SuspiciousLinePattern)
		if err != nil {
			return gtbox_pdf.MarkdownOptions{}, gtbox_pdf.MarkdownReportOptions{}, fmt.Errorf("suspicious_line_pattern: %w", err)
		}
		reportOpts.SuspiciousLinePattern = re
	}

	return mdOpts, reportOpts, nil
}

// trimAnchor strips a leading `^` so the pattern can be re-anchored after the
// `## ` prefix the converter emits.
func trimAnchor(p string) string {
	if strings.HasPrefix(p, "^") {
		return p[1:]
	}
	return p
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
