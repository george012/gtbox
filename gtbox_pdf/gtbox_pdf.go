/*
Package gtbox_pdf provides small PDF text extraction and Markdown cleanup
helpers. It intentionally delegates PDF layout extraction to pdftotext and
keeps the Go side focused on deterministic text cleanup and validation.
*/
package gtbox_pdf

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

type MarkdownProfile string

const (
	ProfileDefault MarkdownProfile = ""
	ProfileNoahAPI MarkdownProfile = "noah-api"
)

type MarkdownOptions struct {
	Profile     MarkdownProfile
	TopHeadings []string
	SubHeadings []string
}

type MarkdownReport struct {
	Lines                int
	CodeFenceCount       int
	EndpointHeadingCount int
	FormFeedCount        int
	ReplacementCharCount int
	SuspiciousURLLines   []int
	NumberedLinesOutside []int
}

var (
	endpointRe     = regexp.MustCompile(`^端点[[:space:]]+[0-9]+:[[:space:]]+POST[[:space:]]+/api/`)
	numberedLineRe = regexp.MustCompile(`^[[:space:]]*[0-9]+[[:space:]]+`)
	loneNumberRe   = regexp.MustCompile(`^[[:space:]]*[0-9]+[[:space:]]*$`)
)

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

func CleanPDFText(raw string) string {
	s := strings.ReplaceAll(raw, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	s = strings.ReplaceAll(s, "\u200b", "")
	s = strings.ReplaceAll(s, "\f", "\n")

	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRightFunc(line, unicode.IsSpace)
	}
	return strings.Join(lines, "\n")
}

func ConvertTextToMarkdown(text string, opts MarkdownOptions) string {
	opts = opts.withDefaults()
	top := stringSet(opts.TopHeadings)
	sub := stringSet(opts.SubHeadings)

	lines := strings.Split(CleanPDFText(text), "\n")
	var out []string
	inCode := false
	pendingCode := false
	skipNoahIndex := false

	emitBlank := func() {
		if len(out) == 0 || out[len(out)-1] != "" {
			out = append(out, "")
		}
	}

	var processLine func(line string)
	processLine = func(line string) {
		trim := strings.TrimSpace(line)

		if skipNoahIndex {
			if trim == "通用协议" {
				skipNoahIndex = false
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
		if isCodeBlockTag(trim) {
			pendingCode = true
			return
		}
		if opts.Profile == ProfileNoahAPI && trim == "服务索引" {
			emitBlank()
			out = append(out, "# 服务索引", "")
			out = append(out, noahServiceIndexMarkdown()...)
			out = append(out, "")
			skipNoahIndex = true
			return
		}
		if endpointRe.MatchString(trim) {
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

func ReportMarkdown(md string) MarkdownReport {
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
		report.ReplacementCharCount += strings.Count(line, "\uFFFD")
		if strings.HasPrefix(trim, "## 端点 ") {
			report.EndpointHeadingCount++
		}
		if !inCode && numberedLineRe.MatchString(line) {
			report.NumberedLinesOutside = append(report.NumberedLinesOutside, lineNo)
		}
		if suspiciousURLFragment(trim) {
			report.SuspiciousURLLines = append(report.SuspiciousURLLines, lineNo)
		}
	}
	return report
}

func (opts MarkdownOptions) withDefaults() MarkdownOptions {
	if len(opts.TopHeadings) == 0 {
		opts.TopHeadings = []string{
			"蛋白质设计函数计算API", "服务索引", "通用协议", "鉴权",
			"Self-description（强烈推荐先调）", "请求格式", "Job 生命周期（异步模式）",
			"JobInfo     数据结构", "JobInfo 数据结构", "Job 查询端点", "轮询建议",
			"错误处理", "跨任务链式调用： job:// URI", "Python 客户端骨架", "文件格式简介",
			"服务部署信息", "限额 / SLA", "RFantibody", "Genie3", "ProteinMPNN",
			"PPIFlow", "PPIflow", "RFdiffusion", "RFDiffusion", "Boltz", "Boltz-2 概念",
			"DockQ", "DockQ 评分速览",
		}
	}
	if len(opts.SubHeadings) == 0 {
		opts.SubHeadings = []string{
			"请求字段", "请求示例", "Python 示例", "输出", "示例", "链式调用示例",
			"FAQ", "Quiver 格式速览", "约束语法", "完整模板", "Schema", "响应",
			"响应示例", "常见错误",
		}
	}
	return opts
}

func stringSet(values []string) map[string]bool {
	m := make(map[string]bool, len(values))
	for _, v := range values {
		m[strings.TrimSpace(v)] = true
	}
	return m
}

func isCodeBlockTag(trim string) bool {
	return trim == "代码块" || trim == "```"
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

func stripLineNumber(line string) string {
	if numberedLineRe.MatchString(line) {
		return numberedLineRe.ReplaceAllString(line, "")
	}
	return strings.TrimLeftFunc(line, unicode.IsSpace)
}

func suspiciousURLFragment(trim string) bool {
	if strings.HasSuffix(trim, "fc.ruosheng.") || strings.HasPrefix(trim, "bio/") {
		return true
	}
	switch trim {
	case "dy", "on", "mpnn":
		return true
	}
	return strings.Contains(trim, "api.colabfo")
}

func noahServiceIndexMarkdown() []string {
	rows := [][]string{
		{"rfantibody", "抗体从头设计三步流水线（RFdiffusion -> ProteinMPNN -> RF2）", "GPU", "https://fc.ruosheng.bio/rfantibody"},
		{"rfdiffusion", "通用蛋白扩散主链生成（unconditional / motif / binder / symmetry / 自定义）", "GPU", "https://fc.ruosheng.bio/rfdiffusion"},
		{"genie3", "通用蛋白扩散生成（unconditional / motif / binder / 自定义 YAML）", "GPU", "https://fc.ruosheng.bio/genie3"},
		{"ppiflow", "PPI flow-matching 结构生成（binder / antibody / nanobody / monomer / scaffolding）", "GPU", "https://fc.ruosheng.bio/ppiflow"},
		{"proteinmpnn", "ProteinMPNN 序列设计 / 打分 / 概率（4 套权重）", "GPU", "https://fc.ruosheng.bio/proteinmpnn"},
		{"boltz", "Boltz-2 复合物结构预测 + ligand 亲和力（AlphaFold3-class）", "GPU", "https://fc.ruosheng.bio/boltz2"},
		{"dockq", "DockQ 复合物结构质量评分（单对 / 批量）", "CPU（待优化）", "https://fc.ruosheng.bio/dockq"},
	}
	out := []string{"| 服务 | 用途 | 算力 | URL |", "|---|---|---|---|"}
	for _, r := range rows {
		out = append(out, fmt.Sprintf("| %s | %s | %s | %s |", r[0], r[1], r[2], r[3]))
	}
	return out
}

func (r MarkdownReport) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "lines=%d\n", r.Lines)
	fmt.Fprintf(&b, "code_fences=%d\n", r.CodeFenceCount)
	fmt.Fprintf(&b, "endpoint_headings=%d\n", r.EndpointHeadingCount)
	fmt.Fprintf(&b, "form_feeds=%d\n", r.FormFeedCount)
	fmt.Fprintf(&b, "replacement_chars=%d\n", r.ReplacementCharCount)
	fmt.Fprintf(&b, "suspicious_url_lines=%s\n", compactInts(r.SuspiciousURLLines, 20))
	fmt.Fprintf(&b, "numbered_lines_outside_code=%s\n", compactInts(r.NumberedLinesOutside, 20))
	return strings.TrimRight(b.String(), "\n")
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
