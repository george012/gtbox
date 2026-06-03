// Package gtbox_cloud_store 云对象存储统一门面。
//
// NewGTCloudStore 按配置中的 provider 分发，返回后端无关的 GTCloudStore 接口；
// 调用方按统一接口读写对象，可在 AWS S3 与阿里云 OSS 之间平滑切换而无需改动业务代码。
//
// 用法：
//
//	store, err := gtbox_cloud_store.NewGTCloudStore(gtbox_cloud_store_common.GTCloudStoreConfig{
//	    Provider:        gtbox_cloud_store_common.GTCloudStoreProvider_AWS_S3,
//	    Endpoint:        "https://s3.us-east-1.amazonaws.com",
//	    Region:          "us-east-1",
//	    AccessKeyID:     "AK...",
//	    AccessKeySecret: "SK...",
//	    Bucket:          "my-bucket",
//	})
//	if err != nil { /* 连接信息缺失会在此显式报错 */ }
//	err = store.PutBytes(ctx, "path/to/object", data, nil)
package gtbox_cloud_store

import (
	"github.com/george012/gtbox/gtbox_cloud_store/gtbox_cloud_store_common"
	"github.com/george012/gtbox/gtbox_cloud_store/gtbox_cloud_store_oss_aliyun"
	"github.com/george012/gtbox/gtbox_cloud_store/gtbox_cloud_store_s3_aws"
)

// NewGTCloudStore 按 provider 构造对应后端的对象存储客户端。
// 未知 provider 返回 ErrGTCloudStoreProviderUnsupported；连接信息校验失败由各实现返回具体错误。
//
// 出错时显式返回 nil 接口：各构造函数返回的是具体指针类型，若直接 return 其 (nil, err)，
// 接口会装箱成「带类型的非 nil 接口」，调用方 store == nil 判断会被误导。此处显式拦截规避该陷阱。
func NewGTCloudStore(cfg gtbox_cloud_store_common.GTCloudStoreConfig) (gtbox_cloud_store_common.GTCloudStore, error) {
	switch cfg.Provider {
	case gtbox_cloud_store_common.GTCloudStoreProvider_AWS_S3:
		store, err := gtbox_cloud_store_s3_aws.NewGTCloudStoreS3AWS(cfg)
		if err != nil {
			return nil, err
		}
		return store, nil
	case gtbox_cloud_store_common.GTCloudStoreProvider_AliYun_Oss:
		store, err := gtbox_cloud_store_oss_aliyun.NewGTCloudStoreOSSAliYun(cfg)
		if err != nil {
			return nil, err
		}
		return store, nil
	default:
		return nil, gtbox_cloud_store_common.ErrGTCloudStoreProviderUnsupported
	}
}
