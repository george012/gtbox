package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/george012/gtbox/gtbox_pdf"
)

func main() {
	in := flag.String("in", "", "input PDF path")
	out := flag.String("out", "", "output Markdown path")
	profile := flag.String("profile", "", "markdown profile: noah-api")
	timeout := flag.Duration("timeout", 2*time.Minute, "pdftotext timeout")
	reportOnly := flag.Bool("report", false, "print report for an existing markdown file")
	flag.Parse()

	if *out == "" {
		fatalf("-out is required")
	}

	if *reportOnly {
		data, err := os.ReadFile(*out)
		if err != nil {
			fatalf("read markdown: %v", err)
		}
		fmt.Println(gtbox_pdf.ReportMarkdown(string(data)).String())
		return
	}

	if *in == "" {
		fatalf("-in is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	text, err := gtbox_pdf.ExtractTextWithPdftotext(ctx, *in)
	if err != nil {
		fatalf("%v", err)
	}

	md := gtbox_pdf.ConvertTextToMarkdown(text, gtbox_pdf.MarkdownOptions{
		Profile: gtbox_pdf.MarkdownProfile(strings.TrimSpace(*profile)),
	})
	if err := os.WriteFile(*out, []byte(md), 0644); err != nil {
		fatalf("write markdown: %v", err)
	}

	fmt.Println(gtbox_pdf.ReportMarkdown(md).String())
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
