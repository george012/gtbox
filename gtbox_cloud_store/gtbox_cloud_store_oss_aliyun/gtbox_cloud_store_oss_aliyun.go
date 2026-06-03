// Package gtbox_cloud_store_oss_aliyun 基于官方 aliyun-oss-go-sdk 实现 GTCloudStore 接口。
//
// 通过 oss.WithContext 把调用方的 context 透传给底层请求，支持超时与取消；
// 元信息从 HEAD 响应头解析，取与 S3 一致的交集字段，保证跨后端可移植。
package gtbox_cloud_store_oss_aliyun

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/george012/gtbox/gtbox_cloud_store/gtbox_cloud_store_common"
)

// GTCloudStoreOSSAliYun 阿里云对象存储 OSS 的客户端。
type GTCloudStoreOSSAliYun struct {
	client     *oss.Client
	bucket     *oss.Bucket
	bucketName string
}

// 编译期断言：实现满足公共契约。
var _ gtbox_cloud_store_common.GTCloudStore = (*GTCloudStoreOSSAliYun)(nil)

// NewGTCloudStoreOSSAliYun 校验配置并构造 OSS 客户端 + 桶句柄。
// 连接信息缺失立即返回错误，不兜底默认 endpoint / bucket。Region 字段 OSS 不使用。
func NewGTCloudStoreOSSAliYun(cfg gtbox_cloud_store_common.GTCloudStoreConfig) (*GTCloudStoreOSSAliYun, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	client, err := oss.New(cfg.Endpoint, cfg.AccessKeyID, cfg.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	bucket, err := client.Bucket(cfg.Bucket)
	if err != nil {
		return nil, err
	}

	return &GTCloudStoreOSSAliYun{
		client:     client,
		bucket:     bucket,
		bucketName: cfg.Bucket,
	}, nil
}

// Provider 返回后端类型。
func (o *GTCloudStoreOSSAliYun) Provider() gtbox_cloud_store_common.GTCloudStoreProvider {
	return gtbox_cloud_store_common.GTCloudStoreProvider_AliYun_Oss
}

// Bucket 返回当前桶名。
func (o *GTCloudStoreOSSAliYun) Bucket() string {
	return o.bucketName
}

// PutObject 流式上传。size >= 0 时显式声明 ContentLength；size < 0 时由 SDK 分块上传。
func (o *GTCloudStoreOSSAliYun) PutObject(ctx context.Context, key string, reader io.Reader, size int64, opts *gtbox_cloud_store_common.GTCloudStorePutOptions) error {
	options := []oss.Option{oss.WithContext(ctx)}
	if size >= 0 {
		options = append(options, oss.ContentLength(size))
	}
	if opts != nil && opts.ContentType != "" {
		options = append(options, oss.ContentType(opts.ContentType))
	}
	return o.bucket.PutObject(key, reader, options...)
}

// PutBytes 上传字节切片。
func (o *GTCloudStoreOSSAliYun) PutBytes(ctx context.Context, key string, data []byte, opts *gtbox_cloud_store_common.GTCloudStorePutOptions) error {
	return o.PutObject(ctx, key, bytes.NewReader(data), int64(len(data)), opts)
}

// GetObject 下载对象，返回的 ReadCloser 由调用方 Close。对象不存在返回 ErrGTCloudStoreObjectNotFound。
func (o *GTCloudStoreOSSAliYun) GetObject(ctx context.Context, key string) (io.ReadCloser, error) {
	rc, err := o.bucket.GetObject(key, oss.WithContext(ctx))
	if err != nil {
		if isOSSNotFound(err) {
			return nil, gtbox_cloud_store_common.ErrGTCloudStoreObjectNotFound
		}
		return nil, err
	}
	return rc, nil
}

// GetBytes 下载对象为字节切片。
func (o *GTCloudStoreOSSAliYun) GetBytes(ctx context.Context, key string) ([]byte, error) {
	rc, err := o.GetObject(ctx, key)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return io.ReadAll(rc)
}

// DeleteObject 删除对象。OSS 删除不存在的键也返回成功，天然幂等。
func (o *GTCloudStoreOSSAliYun) DeleteObject(ctx context.Context, key string) error {
	return o.bucket.DeleteObject(key, oss.WithContext(ctx))
}

// ObjectExists 判断对象是否存在。
func (o *GTCloudStoreOSSAliYun) ObjectExists(ctx context.Context, key string) (bool, error) {
	return o.bucket.IsObjectExist(key, oss.WithContext(ctx))
}

// StatObject 取对象元信息；不存在返回 ErrGTCloudStoreObjectNotFound。
func (o *GTCloudStoreOSSAliYun) StatObject(ctx context.Context, key string) (*gtbox_cloud_store_common.GTCloudStoreObjectInfo, error) {
	header, err := o.bucket.GetObjectDetailedMeta(key, oss.WithContext(ctx))
	if err != nil {
		if isOSSNotFound(err) {
			return nil, gtbox_cloud_store_common.ErrGTCloudStoreObjectNotFound
		}
		return nil, err
	}

	info := &gtbox_cloud_store_common.GTCloudStoreObjectInfo{
		Key:         key,
		ETag:        strings.Trim(header.Get("Etag"), `"`),
		ContentType: header.Get("Content-Type"),
	}
	if v := header.Get("Content-Length"); v != "" {
		if n, perr := strconv.ParseInt(v, 10, 64); perr == nil {
			info.Size = n
		}
	}
	if v := header.Get("Last-Modified"); v != "" {
		if t, perr := http.ParseTime(v); perr == nil {
			info.LastModified = t
		}
	}
	return info, nil
}

// PresignGetURL 生成限时下载预签名 URL。expire 向下取整到秒（OSS 签名按秒）。
func (o *GTCloudStoreOSSAliYun) PresignGetURL(ctx context.Context, key string, expire time.Duration) (string, error) {
	return o.bucket.SignURL(key, oss.HTTPGet, int64(expire/time.Second), oss.WithContext(ctx))
}

// isOSSNotFound 判定错误是否为对象不存在（OSS 返回 oss.ServiceError，键不存在为 HTTP 404 / NoSuchKey）。
func isOSSNotFound(err error) bool {
	var serr oss.ServiceError
	if errors.As(err, &serr) {
		return serr.StatusCode == http.StatusNotFound || serr.Code == "NoSuchKey"
	}
	return false
}
