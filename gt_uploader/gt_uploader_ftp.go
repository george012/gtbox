// Package gt_uploader gt_uploader/gt_uploader_ftp.go
package gt_uploader

import (
	"fmt"
	"io"
)

func uploadFTP(reader io.Reader, filename string, opt UploadOption) (string, error) {
	return "", fmt.Errorf("FTP upload not implemented")
}
