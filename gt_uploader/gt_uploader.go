// Package gt_uploader gt_uploader/gt_uploader.go
package gt_uploader

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type UploadProtocol string

const (
	UploadProtocolHTTP UploadProtocol = "http"
	UploadProtocolSCP  UploadProtocol = "scp"
	UploadProtocolSFTP UploadProtocol = "sftp"
	UploadProtocolFTP  UploadProtocol = "ftp"
)

type UploadOption struct {
	Protocol   UploadProtocol
	Addr       string
	Username   string
	Password   string
	PrivateKey []byte
	Headers    map[string]string
	TargetPath string
	Timeout    time.Duration
}

func Upload(reader io.Reader, filename string, opt UploadOption) (string, error) {
	switch opt.Protocol {
	case UploadProtocolHTTP:
		return uploadHTTP(reader, filename, opt)
	case UploadProtocolSFTP:
		return uploadSFTP(reader, filename, opt)
	case UploadProtocolSCP:
		return uploadSCP(reader, filename, opt)
	case UploadProtocolFTP:
		return uploadFTP(reader, filename, opt)
	default:
		return "", fmt.Errorf("unsupported protocol: %s", opt.Protocol)
	}
}

func UploadDir(localDir string, opt UploadOption, remoteDir string) error {
	info, err := os.Stat(localDir)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("invalid directory: %s", localDir)
	}

	if remoteDir == "" {
		remoteDir = filepath.Base(localDir)
	}

	return filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, _ := filepath.Rel(localDir, path)
		targetPath := filepath.Join(remoteDir, relPath)

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = Upload(file, info.Name(), UploadOption{
			Protocol:   opt.Protocol,
			Addr:       opt.Addr,
			Username:   opt.Username,
			Password:   opt.Password,
			PrivateKey: opt.PrivateKey,
			Headers:    opt.Headers,
			TargetPath: targetPath,
			Timeout:    opt.Timeout,
		})

		return err
	})
}
