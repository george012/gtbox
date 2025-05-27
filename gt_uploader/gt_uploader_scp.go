// Package gt_uploader gt_uploader/gt_uploader_scp.go
package gt_uploader

import (
	"fmt"
	"io"
)

func uploadSCP(reader io.Reader, filename string, opt UploadOption) (string, error) {
	return "", fmt.Errorf("SCP upload not implemented")
}
