// Package gtbox_cloud_store_s3_aws 基于官方 aws-sdk-go-v2 实现 GTCloudStore 接口。
//
// 支持自定义 BaseEndpoint 与 path-style 寻址，因此既可对接真实 AWS S3，也可对接
// 任意 S3 兼容端点（MinIO / localstack / 自建网关）。只用 service/s3 与 credentials
// 两个 SDK 模块、静态凭证手工组装 aws.Config，不引入重量级 config 加载层。
package gtbox_cloud_store_s3_aws

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/george012/gtbox/gtbox_cloud_store/gtbox_cloud_store_common"
)

// GTCloudStoreS3AWS AWS S3（及 S3 兼容端点）的对象存储客户端。
type GTCloudStoreS3AWS struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	bucket        string
}

// 编译期断言：实现满足公共契约。
var _ gtbox_cloud_store_common.GTCloudStore = (*GTCloudStoreS3AWS)(nil)

// NewGTCloudStoreS3AWS 校验配置并构造 S3 客户端。
// 连接信息缺失立即返回错误，不兜底默认 endpoint / region / bucket。
func NewGTCloudStoreS3AWS(cfg gtbox_cloud_store_common.GTCloudStoreConfig) (*GTCloudStoreS3AWS, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	// S3 签名 v4 必须有 region；真实 AWS 用真实地域，自定义端点可填任意非空占位地域。
	if cfg.Region == "" {
		return nil, gtbox_cloud_store_common.ErrGTCloudStoreRegionRequired
	}

	awsCfg := aws.Config{
		Region:      cfg.Region,
		Credentials: credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.AccessKeySecret, ""),
	}
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.Endpoint)
		o.UsePathStyle = cfg.UsePathStyle
	})

	return &GTCloudStoreS3AWS{
		client:        client,
		presignClient: s3.NewPresignClient(client),
		bucket:        cfg.Bucket,
	}, nil
}

// Provider 返回后端类型。
func (s *GTCloudStoreS3AWS) Provider() gtbox_cloud_store_common.GTCloudStoreProvider {
	return gtbox_cloud_store_common.GTCloudStoreProvider_AWS_S3
}

// Bucket 返回当前桶名。
func (s *GTCloudStoreS3AWS) Bucket() string {
	return s.bucket
}

// PutObject 流式上传。size >= 0 时设置 ContentLength；size < 0 时 SDK 会先读满 body
// 计算长度（非 seekable 会进内存），大对象未知长度场景应改走多段上传扩展。
func (s *GTCloudStoreS3AWS) PutObject(ctx context.Context, key string, reader io.Reader, size int64, opts *gtbox_cloud_store_common.GTCloudStorePutOptions) error {
	in := &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   reader,
	}
	if size >= 0 {
		in.ContentLength = aws.Int64(size)
	}
	if opts != nil && opts.ContentType != "" {
		in.ContentType = aws.String(opts.ContentType)
	}
	_, err := s.client.PutObject(ctx, in)
	return err
}

// PutBytes 上传字节切片。
func (s *GTCloudStoreS3AWS) PutBytes(ctx context.Context, key string, data []byte, opts *gtbox_cloud_store_common.GTCloudStorePutOptions) error {
	return s.PutObject(ctx, key, bytes.NewReader(data), int64(len(data)), opts)
}

// GetObject 下载对象，返回的 ReadCloser 由调用方 Close。对象不存在返回 ErrGTCloudStoreObjectNotFound。
func (s *GTCloudStoreS3AWS) GetObject(ctx context.Context, key string) (io.ReadCloser, error) {
	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if isS3NotFound(err) {
			return nil, gtbox_cloud_store_common.ErrGTCloudStoreObjectNotFound
		}
		return nil, err
	}
	return out.Body, nil
}

// GetBytes 下载对象为字节切片。
func (s *GTCloudStoreS3AWS) GetBytes(ctx context.Context, key string) ([]byte, error) {
	rc, err := s.GetObject(ctx, key)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return io.ReadAll(rc)
}

// DeleteObject 删除对象。S3 删除不存在的键也返回成功，天然幂等。
func (s *GTCloudStoreS3AWS) DeleteObject(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}

// ObjectExists 通过 HeadObject 判断对象是否存在。
func (s *GTCloudStoreS3AWS) ObjectExists(ctx context.Context, key string) (bool, error) {
	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if isS3NotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// StatObject 取对象元信息；不存在返回 ErrGTCloudStoreObjectNotFound。
func (s *GTCloudStoreS3AWS) StatObject(ctx context.Context, key string) (*gtbox_cloud_store_common.GTCloudStoreObjectInfo, error) {
	out, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if isS3NotFound(err) {
			return nil, gtbox_cloud_store_common.ErrGTCloudStoreObjectNotFound
		}
		return nil, err
	}

	info := &gtbox_cloud_store_common.GTCloudStoreObjectInfo{
		Key:  key,
		ETag: strings.Trim(aws.ToString(out.ETag), `"`),
	}
	if out.ContentLength != nil {
		info.Size = *out.ContentLength
	}
	if out.ContentType != nil {
		info.ContentType = *out.ContentType
	}
	if out.LastModified != nil {
		info.LastModified = *out.LastModified
	}
	return info, nil
}

// PresignGetURL 生成限时下载预签名 URL。
func (s *GTCloudStoreS3AWS) PresignGetURL(ctx context.Context, key string, expire time.Duration) (string, error) {
	req, err := s.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expire))
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

// isS3NotFound 判定错误是否为对象/键不存在（HeadObject 返回 *types.NotFound，
// GetObject 返回 *types.NoSuchKey，部分兼容端点只回 HTTP 404 错误码）。
func isS3NotFound(err error) bool {
	var notFound *types.NotFound
	if errors.As(err, &notFound) {
		return true
	}
	var noSuchKey *types.NoSuchKey
	if errors.As(err, &noSuchKey) {
		return true
	}
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.ErrorCode() {
		case "NotFound", "NoSuchKey", "404":
			return true
		}
	}
	return false
}
