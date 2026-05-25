package gtbox_pdf

import (
	"strings"
	"testing"
)

func TestConvertTextToMarkdownNoahProfile(t *testing.T) {
	raw := "蛋白质设计函数计算API\n\n服务索引\nrfantibody GPU fc.ruosheng.\nbio/rfantibo\ndy\n通用协议\n\nPPIFlow\n\n端点 1: POST /api/sample/binder\n请求字段\n代码块\n  1   curl -X POST $BASE_URL/api/sample/binder \\\n  2       -F target=@target.pdb\n\nFAQ\nQ: hello\nA: world\n"

	md := ConvertTextToMarkdown(raw, MarkdownOptions{Profile: ProfileNoahAPI})
	for _, want := range []string{
		"# 蛋白质设计函数计算API",
		"# 服务索引",
		"https://fc.ruosheng.bio/rfantibody",
		"# 通用协议",
		"# PPIFlow",
		"## 端点 1: POST /api/sample/binder",
		"### 请求字段",
		"curl -X POST $BASE_URL/api/sample/binder \\",
		"**Q:** hello",
		"**A:** world",
	} {
		if !strings.Contains(md, want) {
			t.Fatalf("markdown missing %q\n%s", want, md)
		}
	}
	report := ReportMarkdown(md)
	if report.FormFeedCount != 0 {
		t.Fatalf("form feed count: got %d", report.FormFeedCount)
	}
	if report.CodeFenceCount%2 != 0 {
		t.Fatalf("unbalanced fences: %d", report.CodeFenceCount)
	}
	if report.EndpointHeadingCount != 1 {
		t.Fatalf("endpoint headings: got %d", report.EndpointHeadingCount)
	}
}

func TestCleanPDFText(t *testing.T) {
	got := CleanPDFText("a\u200b \r\nb\f c\t \n")
	want := "a\nb\n c\n"
	if got != want {
		t.Fatalf("CleanPDFText got %q want %q", got, want)
	}
}
