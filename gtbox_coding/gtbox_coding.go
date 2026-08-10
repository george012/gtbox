package gtbox_coding

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// GTCodeLinesOptions 代码行统计忽略配置
// 忽略规则匹配语义（IgnoreDirs / IgnoreFiles 通用）：
//   - 条目不含 "/"：按名字匹配任意层级，支持 glob（如 "vendor"、"*.pb.go"、"mock_*"）
//   - 条目含 "/"：按项目根相对路径精确匹配（如 "third_party/protos"、"cmd/demo/main.go"）
type GTCodeLinesOptions struct {
	IgnoreDirs      []string // 忽略目录，整棵子树不统计
	IgnoreFiles     []string // 忽略文件
	IgnoreTestFiles bool     // true 时跳过 *_test.go
}

// GetProjectCodeLines 获取当前项目的有效代码行数（全部 *.go，不忽略任何文件）
func GetProjectCodeLines() int64 {
	projectRoot := projectRootFromCaller()
	if projectRoot == "" {
		return 0
	}
	return countGoFileLines(projectRoot, nil)
}

// GetProjectCodeLinesWithOptions 按忽略配置获取当前项目的有效代码行数
func GetProjectCodeLinesWithOptions(options *GTCodeLinesOptions) int64 {
	projectRoot := projectRootFromCaller()
	if projectRoot == "" {
		return 0
	}
	return countGoFileLines(projectRoot, options)
}

// projectRootFromCaller 从调用 GetProjectCodeLines* 的文件位置向上查找 go.mod 所在目录
// 注意：Caller(2) 依赖固定调用深度，本函数只能被本包导出函数直接调用
func projectRootFromCaller() string {
	_, filename, _, ok := runtime.Caller(2)
	if !ok {
		return ""
	}
	currentDir := filepath.Dir(filename)
	for {
		if _, err := os.Stat(filepath.Join(currentDir, "go.mod")); err == nil {
			return currentDir
		}
		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			return ""
		}
		currentDir = parentDir
	}
}

// countGoFileLines 遍历项目根下全部 *.go 文件累计行数（与 wc -l 同语义，按 '\n' 计数）
func countGoFileLines(projectRoot string, options *GTCodeLinesOptions) int64 {
	var total int64
	_ = filepath.WalkDir(projectRoot, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			if entry != nil && entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		relPath, relErr := filepath.Rel(projectRoot, path)
		if relErr != nil {
			return nil
		}
		relPath = filepath.ToSlash(relPath)
		if entry.IsDir() {
			if relPath != "." && options != nil && matchIgnoreRule(options.IgnoreDirs, relPath, entry.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		fileName := entry.Name()
		if !strings.HasSuffix(fileName, ".go") {
			return nil
		}
		if options != nil {
			if options.IgnoreTestFiles && strings.HasSuffix(fileName, "_test.go") {
				return nil
			}
			if matchIgnoreRule(options.IgnoreFiles, relPath, fileName) {
				return nil
			}
		}
		total += fileLineCount(path)
		return nil
	})
	return total
}

// matchIgnoreRule 见 GTCodeLinesOptions 注释的匹配语义
func matchIgnoreRule(rules []string, relPath string, baseName string) bool {
	for _, rule := range rules {
		rule = filepath.ToSlash(strings.TrimSpace(rule))
		rule = strings.TrimPrefix(rule, "./")
		rule = strings.TrimSuffix(rule, "/")
		if rule == "" {
			continue
		}
		if strings.Contains(rule, "/") {
			if relPath == rule || strings.HasPrefix(relPath, rule+"/") {
				return true
			}
			continue
		}
		if matched, _ := filepath.Match(rule, baseName); matched {
			return true
		}
	}
	return false
}

// fileLineCount 单文件行数，读失败按 0 行跳过不影响整体统计
func fileLineCount(filePath string) int64 {
	f, err := os.Open(filePath)
	if err != nil {
		return 0
	}
	defer f.Close()
	var count int64
	buf := make([]byte, 64*1024)
	for {
		n, readErr := f.Read(buf)
		count += int64(bytes.Count(buf[:n], []byte{'\n'}))
		if readErr != nil {
			return count
		}
	}
}
