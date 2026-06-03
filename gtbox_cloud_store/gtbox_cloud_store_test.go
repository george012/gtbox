package gtbox_cloud_store

import (
	"errors"
	"testing"

	"github.com/george012/gtbox/gtbox_cloud_store/gtbox_cloud_store_common"
)

// 一份各字段齐全的 S3 基准配置，单测里按需置空某字段以验证校验逻辑。
func fullS3Config() gtbox_cloud_store_common.GTCloudStoreConfig {
	return gtbox_cloud_store_common.GTCloudStoreConfig{
		Provider:        gtbox_cloud_store_common.GTCloudStoreProvider_AWS_S3,
		Endpoint:        "https://s3.us-east-1.amazonaws.com",
		Region:          "us-east-1",
		AccessKeyID:     "AKIATEST",
		AccessKeySecret: "secret",
		Bucket:          "test-bucket",
	}
}

// TestNewGTCloudStore_UnsupportedProvider 未知 / 空 provider 必须显式拒绝，不静默返回 nil。
func TestNewGTCloudStore_UnsupportedProvider(t *testing.T) {
	for _, p := range []gtbox_cloud_store_common.GTCloudStoreProvider{"", "azure", "gcs"} {
		cfg := fullS3Config()
		cfg.Provider = p
		store, err := NewGTCloudStore(cfg)
		if !errors.Is(err, gtbox_cloud_store_common.ErrGTCloudStoreProviderUnsupported) {
			t.Fatalf("provider=%q want ErrGTCloudStoreProviderUnsupported, got err=%v", p, err)
		}
		if store != nil {
			t.Fatalf("provider=%q want nil store on error, got %v", p, store)
		}
	}
}

// TestNewGTCloudStore_MissingConnInfo 连接信息任一缺失都必须在构造时立即报对应错误，绝不兜底。
func TestNewGTCloudStore_MissingConnInfo(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*gtbox_cloud_store_common.GTCloudStoreConfig)
		wantErr error
	}{
		{"缺 endpoint", func(c *gtbox_cloud_store_common.GTCloudStoreConfig) { c.Endpoint = "" }, gtbox_cloud_store_common.ErrGTCloudStoreEndpointRequired},
		{"缺 access key id", func(c *gtbox_cloud_store_common.GTCloudStoreConfig) { c.AccessKeyID = "" }, gtbox_cloud_store_common.ErrGTCloudStoreAccessKeyIDRequired},
		{"缺 access key secret", func(c *gtbox_cloud_store_common.GTCloudStoreConfig) { c.AccessKeySecret = "" }, gtbox_cloud_store_common.ErrGTCloudStoreAccessKeySecretRequired},
		{"缺 bucket", func(c *gtbox_cloud_store_common.GTCloudStoreConfig) { c.Bucket = "" }, gtbox_cloud_store_common.ErrGTCloudStoreBucketRequired},
		{"S3 缺 region", func(c *gtbox_cloud_store_common.GTCloudStoreConfig) { c.Region = "" }, gtbox_cloud_store_common.ErrGTCloudStoreRegionRequired},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := fullS3Config()
			tt.mutate(&cfg)
			store, err := NewGTCloudStore(cfg)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("want %v, got %v", tt.wantErr, err)
			}
			if store != nil {
				t.Fatalf("want nil store on error, got %v", store)
			}
		})
	}
}

// TestNewGTCloudStore_S3OK 完整 S3 配置应构造成功，且句柄报告正确的后端与桶名（构造不发起网络）。
func TestNewGTCloudStore_S3OK(t *testing.T) {
	store, err := NewGTCloudStore(fullS3Config())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := store.Provider(); got != gtbox_cloud_store_common.GTCloudStoreProvider_AWS_S3 {
		t.Fatalf("provider want aws_s3, got %q", got)
	}
	if got := store.Bucket(); got != "test-bucket" {
		t.Fatalf("bucket want test-bucket, got %q", got)
	}
}

// TestNewGTCloudStore_OSSOK 完整 OSS 配置应构造成功；OSS 不使用 Region，缺省也不应报错。
func TestNewGTCloudStore_OSSOK(t *testing.T) {
	cfg := fullS3Config()
	cfg.Provider = gtbox_cloud_store_common.GTCloudStoreProvider_AliYun_Oss
	cfg.Endpoint = "https://oss-cn-hangzhou.aliyuncs.com"
	cfg.Region = "" // OSS 忽略 region
	store, err := NewGTCloudStore(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := store.Provider(); got != gtbox_cloud_store_common.GTCloudStoreProvider_AliYun_Oss {
		t.Fatalf("provider want aliyun_oss, got %q", got)
	}
	if got := store.Bucket(); got != "test-bucket" {
		t.Fatalf("bucket want test-bucket, got %q", got)
	}
}
