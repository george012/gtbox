package gtbox_pdf

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

// MarkdownOptions configures ConvertTextToMarkdown. Every list / pattern is
// optional and supplied by the caller; the package itself ships no
// domain-specific rules.
type MarkdownOptions struct {
	// Lines whose trimmed text exactly matches one of TopHeadings become
	// "# X" in the output.
	TopHeadings []string

	// Lines whose trimmed text exactly matches one of SubHeadings become
	// "### X" in the output.
	SubHeadings []string

	// Lines whose trimmed text matches SectionHeadingPattern become
	// "## <trimmed-line>" in the output. nil disables this transformation.
	SectionHeadingPattern *regexp.Regexp

	// CodeBlockMarkers lists trimmed line values that switch the converter
	// into "pending code" mode: the next non-empty line opens a fenced
	// block. When nil, ["```"] is used. Pass an empty (non-nil) slice to
	// disable the marker mechanism entirely.
	CodeBlockMarkers []string

	// LineJoinRules describe multi-line repairs. They run after CleanPDFText
	// and before structural rewrites.
	LineJoinRules []MarkdownLineJoinRule

	// Injections inject pre-composed Markdown blocks when a trigger line is
	// found. Useful when an upstream PDF degrades a section (e.g. a service
	// index table) and the caller wants to substitute a clean version.
	Injections []MarkdownInjection

	// QAEnabled, when true, rewrites lines starting with "Q:" / "A:" into
	// "**Q:** ..." / "**A:** ...". Default false.
	QAEnabled bool
}

// MarkdownLineJoinRule describes a two-line repair. When line i contains
// Contains and the trimmed line i+1 equals NextEquals, line i has Replace
// substituted with With and line i+1 is dropped.
type MarkdownLineJoinRule struct {
	Contains   string
	NextEquals string
	Replace    string
	With       string
}

// MarkdownInjection replaces a region of the input with a caller-supplied
// Markdown block.
type MarkdownInjection struct {
	// Trigger is the trimmed line value that fires the injection.
	Trigger string
	// Heading, when non-empty, is emitted first as "# Heading".
	Heading string
	// Lines are the raw Markdown lines emitted after the heading.
	Lines []string
	// SkipUntil, when non-empty, causes subsequent input lines to be
	// dropped until a line whose trimmed text equals SkipUntil appears;
	// that match line is then re-processed normally.
	SkipUntil string
}

// MarkdownReport captures structural counters for a Markdown document.
type MarkdownReport struct {
	Lines                int
	CodeFenceCount       int
	SectionHeadingCount  int
	FormFeedCount        int
	ReplacementCharCount int
	SuspiciousLines      []int
	NumberedLinesOutside []int
}

// MarkdownReportOptions configures pattern-driven counters for
// ReportMarkdown.
type MarkdownReportOptions struct {
	// SectionHeadingPattern counts lines whose trimmed text matches it.
	// Defaults to `^## ` when nil.
	SectionHeadingPattern *regexp.Regexp
	// SuspiciousLinePattern, when non-nil, records line numbers of lines
	// whose trimmed text matches.
	SuspiciousLinePattern *regexp.Regexp
}

var (
	defaultCodeBlockMarkers = []string{"```"}
	defaultSectionHeadingRe = regexp.MustCompile(`^## `)
)

// ConvertTextToMarkdown rewrites extracted PDF text into Markdown using only
// the rules supplied by opts. The function is deterministic and
// side-effect-free.
func ConvertTextToMarkdown(text string, opts MarkdownOptions) string {
	top := stringSet(opts.TopHeadings)
	sub := stringSet(opts.SubHeadings)
	fence := opts.CodeBlockMarkers
	if fence == nil {
		fence = defaultCodeBlockMarkers
	}
	fenceSet := stringSet(fence)
	injectionsByTrigger := indexInjections(opts.Injections)

	lines := applyLineJoinRules(strings.Split(CleanPDFText(text), "\n"), opts.LineJoinRules)
	var out []string
	inCode := false
	pendingCode := false
	skipUntil := ""

	emitBlank := func() {
		if len(out) == 0 || out[len(out)-1] != "" {
			out = append(out, "")
		}
	}

	var processLine func(line string)
	processLine = func(line string) {
		trim := strings.TrimSpace(line)

		if skipUntil != "" {
			if trim == skipUntil {
				skipUntil = ""
				processLine(line)
			}
			return
		}

		if pendingCode {
			if trim == "" || loneNumberRe.MatchString(line) {
				return
			}
			if isCodeLineCandidate(line) {
				emitBlank()
				out = append(out, "```", stripLineNumber(line))
				inCode = true
				pendingCode = false
				return
			}
			pendingCode = false
		}

		if inCode {
			if trim == "" {
				out = append(out, "")
				return
			}
			if isCodeLineCandidate(line) {
				out = append(out, stripLineNumber(line))
				return
			}
			out = append(out, "```", "")
			inCode = false
		}

		if trim == "" || loneNumberRe.MatchString(line) {
			emitBlank()
			return
		}
		if fenceSet[trim] {
			pendingCode = true
			return
		}
		if numberedLineRe.MatchString(line) && isLikelyCodeContent(stripLineNumber(line)) {
			emitBlank()
			out = append(out, "```", stripLineNumber(line))
			inCode = true
			return
		}
		if inj, ok := injectionsByTrigger[trim]; ok {
			emitBlank()
			if inj.Heading != "" {
				out = append(out, "# "+inj.Heading, "")
			}
			out = append(out, inj.Lines...)
			out = append(out, "")
			if inj.SkipUntil != "" {
				skipUntil = inj.SkipUntil
			}
			return
		}
		if opts.SectionHeadingPattern != nil && opts.SectionHeadingPattern.MatchString(trim) {
			emitBlank()
			out = append(out, "## "+trim)
			return
		}
		if top[trim] {
			emitBlank()
			out = append(out, "# "+trim)
			return
		}
		if sub[trim] {
			emitBlank()
			out = append(out, "### "+trim)
			return
		}
		if opts.QAEnabled {
			if strings.HasPrefix(trim, "Q:") {
				emitBlank()
				out = append(out, "**Q:** "+strings.TrimSpace(strings.TrimPrefix(trim, "Q:")))
				return
			}
			if strings.HasPrefix(trim, "A:") {
				emitBlank()
				out = append(out, "**A:** "+strings.TrimSpace(strings.TrimPrefix(trim, "A:")))
				return
			}
		}
		out = append(out, line)
	}

	for _, line := range lines {
		processLine(line)
	}
	if inCode {
		out = append(out, "```")
	}
	return strings.TrimSpace(strings.Join(out, "\n")) + "\n"
}

// ReportMarkdown produces structural counters for a Markdown document.
func ReportMarkdown(md string, opts MarkdownReportOptions) MarkdownReport {
	section := opts.SectionHeadingPattern
	if section == nil {
		section = defaultSectionHeadingRe
	}
	lines := strings.Split(md, "\n")
	report := MarkdownReport{Lines: len(lines)}
	inCode := false
	for i, line := range lines {
		lineNo := i + 1
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "```") {
			report.CodeFenceCount++
			inCode = !inCode
			continue
		}
		report.FormFeedCount += strings.Count(line, "\f")
		report.ReplacementCharCount += strings.Count(line, "�")
		if section.MatchString(trim) {
			report.SectionHeadingCount++
		}
		if !inCode && numberedLineRe.MatchString(line) {
			report.NumberedLinesOutside = append(report.NumberedLinesOutside, lineNo)
		}
		if opts.SuspiciousLinePattern != nil && opts.SuspiciousLinePattern.MatchString(trim) {
			report.SuspiciousLines = append(report.SuspiciousLines, lineNo)
		}
	}
	return report
}

// String renders the report as a stable key=value text block.
func (r MarkdownReport) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "lines=%d\n", r.Lines)
	fmt.Fprintf(&b, "code_fences=%d\n", r.CodeFenceCount)
	fmt.Fprintf(&b, "section_headings=%d\n", r.SectionHeadingCount)
	fmt.Fprintf(&b, "form_feeds=%d\n", r.FormFeedCount)
	fmt.Fprintf(&b, "replacement_chars=%d\n", r.ReplacementCharCount)
	fmt.Fprintf(&b, "suspicious_lines=%s\n", compactInts(r.SuspiciousLines, 20))
	fmt.Fprintf(&b, "numbered_lines_outside_code=%s\n", compactInts(r.NumberedLinesOutside, 20))
	return strings.TrimRight(b.String(), "\n")
}

func stringSet(values []string) map[string]bool {
	m := make(map[string]bool, len(values))
	for _, v := range values {
		m[strings.TrimSpace(v)] = true
	}
	return m
}

func indexInjections(injections []MarkdownInjection) map[string]MarkdownInjection {
	if len(injections) == 0 {
		return nil
	}
	m := make(map[string]MarkdownInjection, len(injections))
	for _, inj := range injections {
		m[strings.TrimSpace(inj.Trigger)] = inj
	}
	return m
}

func isCodeLineCandidate(line string) bool {
	if numberedLineRe.MatchString(line) {
		return true
	}
	if len(line) == 0 {
		return false
	}
	return unicode.IsSpace([]rune(line)[0])
}

func isLikelyCodeContent(line string) bool {
	trim := strings.TrimSpace(line)
	if trim == "" {
		return false
	}
	codePrefixes := []string{
		"#", "import ", "from ", "def ", "class ", "with ", "for ", "while ", "if ", "elif ",
		"else", "return ", "raise ", "assert ", "print(", "curl ", "BASE_URL", "TOKEN",
	}
	for _, p := range codePrefixes {
		if strings.HasPrefix(trim, p) {
			return true
		}
	}
	codeContains := []string{
		" = ", ":=", ".post(", ".get(", ".json()", ".raise_for_status()", ".iter_bytes()",
		"open(", "Path(", "time.", "httpx.", "f\"",
	}
	for _, part := range codeContains {
		if strings.Contains(trim, part) {
			return true
		}
	}
	return strings.HasSuffix(trim, "{") || strings.HasSuffix(trim, "\\") ||
		strings.HasSuffix(trim, ",") || strings.HasSuffix(trim, ")") ||
		strings.HasSuffix(trim, "]") || strings.HasSuffix(trim, "}")
}

func stripLineNumber(line string) string {
	if numberedLineRe.MatchString(line) {
		return numberedLineRe.ReplaceAllString(line, "")
	}
	return strings.TrimLeftFunc(line, unicode.IsSpace)
}

func applyLineJoinRules(lines []string, rules []MarkdownLineJoinRule) []string {
	if len(rules) == 0 {
		return lines
	}
	repaired := make([]string, 0, len(lines))
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		matched := false
		if i+1 < len(lines) {
			next := strings.TrimSpace(lines[i+1])
			for _, rule := range rules {
				if rule.NextEquals == next && strings.Contains(line, rule.Contains) {
					repaired = append(repaired, strings.Replace(line, rule.Replace, rule.With, 1))
					i++
					matched = true
					break
				}
			}
		}
		if !matched {
			repaired = append(repaired, line)
		}
	}
	return repaired
}

func compactInts(values []int, limit int) string {
	if len(values) == 0 {
		return "[]"
	}
	cp := append([]int(nil), values...)
	sort.Ints(cp)
	if len(cp) > limit {
		cp = cp[:limit]
		return fmt.Sprintf("%v...(+%d)", cp, len(values)-limit)
	}
	return fmt.Sprintf("%v", cp)
}
