package gtbox_coding

import (
	"os"
	"path/filepath"
	"testing"
)

// buildTestProject 构造临时项目树，返回项目根
// 结构：
//
//	main.go            3 行
//	main_test.go       5 行
//	vendor/dep.go      7 行
//	api/api.pb.go      11 行
//	api/api.go         2 行
//	third_party/x/a.go 13 行
func buildTestProject(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	files := map[string]int{
		"main.go":            3,
		"main_test.go":       5,
		"vendor/dep.go":      7,
		"api/api.pb.go":      11,
		"api/api.go":         2,
		"third_party/x/a.go": 13,
	}
	for relPath, lines := range files {
		fullPath := filepath.Join(root, filepath.FromSlash(relPath))
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			t.Fatal(err)
		}
		content := make([]byte, 0, lines*10)
		for range lines {
			content = append(content, []byte("// line\n")...)
		}
		if err := os.WriteFile(fullPath, content, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestCountGoFileLines(t *testing.T) {
	root := buildTestProject(t)

	cases := []struct {
		name    string
		options *GTCodeLinesOptions
		want    int64
	}{
		{"无忽略配置全量统计", nil, 3 + 5 + 7 + 11 + 2 + 13},
		{"忽略测试文件", &GTCodeLinesOptions{IgnoreTestFiles: true}, 3 + 7 + 11 + 2 + 13},
		{"按目录名忽略任意层级", &GTCodeLinesOptions{IgnoreDirs: []string{"vendor"}}, 3 + 5 + 11 + 2 + 13},
		{"按相对路径忽略子树", &GTCodeLinesOptions{IgnoreDirs: []string{"third_party/x"}}, 3 + 5 + 7 + 11 + 2},
		{"按 glob 忽略生成文件", &GTCodeLinesOptions{IgnoreFiles: []string{"*.pb.go"}}, 3 + 5 + 7 + 2 + 13},
		{"按相对路径忽略单文件", &GTCodeLinesOptions{IgnoreFiles: []string{"api/api.go"}}, 3 + 5 + 7 + 11 + 13},
		{"组合配置", &GTCodeLinesOptions{
			IgnoreDirs:      []string{"vendor", "third_party"},
			IgnoreFiles:     []string{"*.pb.go"},
			IgnoreTestFiles: true,
		}, 3 + 2},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := countGoFileLines(root, tc.options); got != tc.want {
				t.Fatalf("got %d, want %d", got, tc.want)
			}
		})
	}
}

func TestMatchIgnoreRule(t *testing.T) {
	cases := []struct {
		rules    []string
		relPath  string
		baseName string
		want     bool
	}{
		{[]string{"vendor"}, "sub/vendor", "vendor", true},        // 裸名匹配任意层级
		{[]string{"vendor/"}, "vendor", "vendor", true},           // 尾部斜杠归一
		{[]string{"./cmd/demo"}, "cmd/demo", "demo", true},        // ./ 前缀归一
		{[]string{"cmd/demo"}, "cmd/demo/sub", "sub", true},       // 路径前缀盖住子树
		{[]string{"cmd/demo"}, "cmd/demofake", "demofake", false}, // 前缀不越界
		{[]string{"*.pb.go"}, "api/api.pb.go", "api.pb.go", true}, // glob 按文件名
		{[]string{"*.pb.go"}, "api/api.go", "api.go", false},      //
		{[]string{"", "  "}, "main.go", "main.go", false},         // 空条目跳过
	}
	for _, tc := range cases {
		if got := matchIgnoreRule(tc.rules, tc.relPath, tc.baseName); got != tc.want {
			t.Fatalf("rules=%v relPath=%q baseName=%q: got %v, want %v", tc.rules, tc.relPath, tc.baseName, got, tc.want)
		}
	}
}

// TestGetProjectCodeLinesAgainstRepo 零参版应能在 gtbox 仓自身统计出非 0 行数
func TestGetProjectCodeLinesAgainstRepo(t *testing.T) {
	if got := GetProjectCodeLines(); got <= 0 {
		t.Fatalf("got %d, want > 0", got)
	}
	withOptions := GetProjectCodeLinesWithOptions(&GTCodeLinesOptions{IgnoreTestFiles: true})
	full := GetProjectCodeLines()
	if withOptions <= 0 || withOptions >= full {
		t.Fatalf("ignore tests got %d, full %d, want 0 < got < full", withOptions, full)
	}
}
