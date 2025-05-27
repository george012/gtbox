package gt_uploader

import (
	"fmt"
	"io"
)

// Placeholder implementations for other protocols
func uploadSFTP(reader io.Reader, filename string, opt UploadOption) (string, error) {
	return "", fmt.Errorf("SFTP upload not implemented")
}
