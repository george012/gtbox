// Package gt_uploader gt_uploader/gt_uploader_http.go
package gt_uploader

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Minimal HTTP uploader (example only)
func uploadHTTP(reader io.Reader, filename string, opt UploadOption) (string, error) {
	body := &strings.Builder{}
	_, err := io.Copy(body, reader)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("PUT", opt.Addr+"/"+opt.TargetPath, strings.NewReader(body.String()))
	if err != nil {
		return "", err
	}
	for k, v := range opt.Headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: opt.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return opt.Addr + "/" + opt.TargetPath, nil
	}
	return "", fmt.Errorf("upload failed with status: %s", resp.Status)
}
