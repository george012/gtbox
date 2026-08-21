// Package gtbox_cloud_store_common 定义云对象存储的可移植契约层。
//
// 这里只放后端无关的东西：统一接口 GTCloudStore、连接配置 GTCloudStoreConfig、
// 对象元信息 GTCloudStoreObjectInfo 与统一错误。各 provider 实现（AWS S3 /
// 阿里云 OSS）只需满足本接口，调用方对具体后端无感、可平滑切换。
package gtbox_cloud_store_common

import (
	"context"
	"errors"
	"io"
	"time"
)

// GTCloudStoreProvider 云存储后端类型。
type GTCloudStoreProvider string

const (
	// GTCloudStoreProvider_AWS_S3 AWS S3 及任意 S3 兼容端点（MinIO / 自建网关）
	GTCloudStoreProvider_AWS_S3 GTCloudStoreProvider = "aws_s3"
	// GTCloudStoreProvider_AliYun_Oss 阿里云对象存储 OSS
	GTCloudStoreProvider_AliYun_Oss GTCloudStoreProvider = "aliyun_oss"
)

// 连接信息缺失 / 对象不存在 / 后端不支持的统一错误。
// 构造各 provider 客户端前做显式校验，零值连接信息一律拒绝，不兜底默认值。
var (
	ErrGTCloudStoreEndpointRequired        = errors.New("gtbox_cloud_store: endpoint is required")
	ErrGTCloudStoreAccessKeyIDRequired     = errors.New("gtbox_cloud_store: access key id is required")
	ErrGTCloudStoreAccessKeySecretRequired = errors.New("gtbox_cloud_store: access key secret is required")
	ErrGTCloudStoreBucketRequired          = errors.New("gtbox_cloud_store: bucket is required")
	ErrGTCloudStoreRegionRequired          = errors.New("gtbox_cloud_store: region is required")
	ErrGTCloudStoreProviderUnsupported     = errors.New("gtbox_cloud_store: provider is unsupported")
	ErrGTCloudStoreObjectNotFound          = errors.New("gtbox_cloud_store: object not found")
)

// GTCloudStoreConfig 云存储连接配置。
//
// 连接信息（Endpoint / AccessKeyID / AccessKeySecret / Bucket）全部必填、显式传入；
// 零值由构造函数立即拒绝，绝不兜底成默认 endpoint / region / bucket —— 那是某个业务的
// 假设，不是公用层应承诺的契约。Region 仅 S3 签名需要，由 S3 实现自行校验。
type GTCloudStoreConfig struct {
	Provider        GTCloudStoreProvider // 后端类型
	Endpoint        string               // 服务端点：S3 可为自定义端点(MinIO/S3兼容)，OSS 为地域端点
	Region          string               // S3 地域，签名需要；OSS 忽略
	AccessKeyID     string               // 访问密钥 ID
	AccessKeySecret string               // 访问密钥 Secret
	Bucket          string               // 桶名
	UsePathStyle    bool                 // S3 path-style 寻址：对接 MinIO/自定义端点时置 true；OSS 忽略
}

// Validate 校验后端无关的必填连接信息。Region 等 provider 专属字段由对应实现校验。
func (cfg GTCloudStoreConfig) Validate() error {
	if cfg.Endpoint == "" {
		return ErrGTCloudStoreEndpointRequired
	}
	if cfg.AccessKeyID == "" {
		return ErrGTCloudStoreAccessKeyIDRequired
	}
	if cfg.AccessKeySecret == "" {
		return ErrGTCloudStoreAccessKeySecretRequired
	}
	if cfg.Bucket == "" {
		return ErrGTCloudStoreBucketRequired
	}
	return nil
}

// GTCloudStoreObjectInfo 对象元信息，取两端交集字段以保证可移植。
type GTCloudStoreObjectInfo struct {
	Key          string    // 对象键
	Size         int64     // 字节大小
	ETag         string    // 内容指纹（已去除两端引号）
	ContentType  string    // 内容类型
	LastModified time.Time // 最后修改时间
}

// GTCloudStorePutOptions 上传可选项。传 nil 表示全部取后端默认。
type GTCloudStorePutOptions struct {
	ContentType string // 内容类型；空则交由后端按默认处理
}

// GTCloudStore 可移植对象存储接口。
//
// 所有方法接受 context.Context 以支持超时与取消。该接口只承诺两端语义一致、
// 字节级可移植的操作；provider 专属能力（多段上传、生命周期、ACL 等）不进本接口，
// 需要时走具体实现的扩展方法，避免把单端特性渗进公共契约。
type GTCloudStore interface {
	// Provider 返回当前后端类型
	Provider() GTCloudStoreProvider
	// Bucket 返回当前操作的桶名
	Bucket() string
	// PutObject 流式上传。size 为内容长度；size < 0 表示长度未知，由实现自行处理
	PutObject(ctx context.Context, key string, reader io.Reader, size int64, opts *GTCloudStorePutOptions) error
	// PutBytes 上传字节切片
	PutBytes(ctx context.Context, key string, data []byte, opts *GTCloudStorePutOptions) error
	// GetObject 下载对象，返回的 ReadCloser 由调用方负责 Close
	GetObject(ctx context.Context, key string) (io.ReadCloser, error)
	// GetBytes 下载对象为字节切片
	GetBytes(ctx context.Context, key string) ([]byte, error)
	// DeleteObject 删除对象。对象不存在视为成功，保证幂等
	DeleteObject(ctx context.Context, key string) error
	// ObjectExists 判断对象是否存在
	ObjectExists(ctx context.Context, key string) (bool, error)
	// StatObject 取对象元信息；不存在返回 ErrGTCloudStoreObjectNotFound
	StatObject(ctx context.Context, key string) (*GTCloudStoreObjectInfo, error)
	// PresignGetURL 生成限时下载预签名 URL
	PresignGetURL(ctx context.Context, key string, expire time.Duration) (string, error)
}
