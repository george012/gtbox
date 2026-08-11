/*
Package gtbox_redis en: Redis handle, zh-cn: Redis封装处理

默认单例 + 多实例并存(对齐 gtbox_log 模式):
  - 单例:SetupRedisConnection 一次初始化,业务侧直接 OwnerRedis.Xxx 使用
  - 多实例:New 创建独立句柄(独立 prefix / 独立连接池),方法集与单例完全一致
*/
package gtbox_redis

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/george012/gtbox/gtbox_log"
	"github.com/redis/go-redis/v9"
)

var (
	// OwnerRedis 默认单例,SetupRedisConnection 后可用
	OwnerRedis *GTRedis
	ownerOnce  sync.Once
	ctx        = context.Background()
)

type RedisConfig struct {
	Addr       string `yaml:"addr" json:"addr"`              // address
	Pwd        string `yaml:"pwd" json:"pwd"`                // pwd
	SocketBuck int    `yaml:"socketBuck" json:"socket_buck"` // 插槽(redis DB index)
	Username   string `yaml:"username" json:"username"`      // ACL 账号,空 = 传统 requirepass 鉴权
	ReadOnly   bool   `yaml:"readOnly" json:"read_only"`     // 句柄级只读防呆,写方法直接拒绝;权限边界靠服务端 ACL
	// PoolSizeMultiplier 默认池容量倍率:连接上限 = 倍率 × 10 × GOMAXPROCS,0 = 默认 1 倍。
	// 池为全进程共享结构,连接不绑核;GOMAXPROCS 仅作机器规格代理,倍率随部署机自适应。
	PoolSizeMultiplier int `yaml:"poolSizeMultiplier" json:"pool_size_multiplier"`
}

type GTRedis struct {
	cfg         *RedisConfig
	prefix      string
	redisClient *redis.Client
}

// New 创建独立 Redis 实例。error 仅参数校验;redis.NewClient 纯内存构造不拨号,
// 真连接首次命令时懒建立,断线由连接池自动重连。
func New(redisCfg RedisConfig, prefixStr string) (*GTRedis, error) {
	if redisCfg.Addr == "" {
		return nil, fmt.Errorf("gtbox_redis: Addr is required")
	}
	if redisCfg.PoolSizeMultiplier < 0 {
		return nil, fmt.Errorf("gtbox_redis: PoolSizeMultiplier must be >= 0, got %d", redisCfg.PoolSizeMultiplier)
	}

	opts := &redis.Options{
		Addr:     redisCfg.Addr,
		Username: redisCfg.Username,
		Password: redisCfg.Pwd,
		DB:       redisCfg.SocketBuck,
	}
	if redisCfg.PoolSizeMultiplier > 0 {
		opts.PoolSize = redisCfg.PoolSizeMultiplier * 10 * runtime.GOMAXPROCS(0)
	}

	return &GTRedis{
		cfg:         &redisCfg,
		prefix:      prefixStr,
		redisClient: redis.NewClient(opts),
	}, nil
}

// SetupRedisConnection 初始化默认单例(sync.Once 收口,并发安全)。
// 返回值为 PING 连通性告知:false 仅表示当前 PING 未通(或参数非法),单例仍已就绪,
// 后续操作由连接池懒拨号 + 自动重连;只读 ACL 账号可正常接入。
func SetupRedisConnection(redisCfg RedisConfig, prefixStr string) (success bool) {
	ownerOnce.Do(func() {
		aRedis, err := New(redisCfg, prefixStr)
		if err != nil {
			gtbox_log.LogErrorf("[redis setup] error [%s]", err)
			return
		}
		OwnerRedis = aRedis
	})
	if OwnerRedis == nil {
		return false
	}
	if err := OwnerRedis.redisClient.Ping(ctx).Err(); err != nil {
		gtbox_log.LogErrorf("[redis setup] ping error [%s]", err)
		return false
	}
	gtbox_log.LogInfof("[redis setup] [%s]", "success")
	return true
}

// Close 关闭连接池,释放全部连接
func (gtr *GTRedis) Close() error {
	return gtr.redisClient.Close()
}

// writeGuard ReadOnly 句柄的写防呆;返回 nil 表示允许写
func (gtr *GTRedis) writeGuard() error {
	if gtr.cfg.ReadOnly {
		return fmt.Errorf("gtbox_redis: write rejected, instance is read-only [prefix=%s]", gtr.prefix)
	}
	return nil
}

// Set 插入单条数据(持久键,无过期;需过期用 SetEX)
func (gtr *GTRedis) Set(key string, value string) error {
	if err := gtr.writeGuard(); err != nil {
		return err
	}
	aKey := fmt.Sprintf("%s:%s", gtr.prefix, key)
	err := gtr.redisClient.Set(ctx, aKey, value, 0).Err()

	return err
}

// SetEX 插入单条数据并设置过期时间
func (gtr *GTRedis) SetEX(key string, value string, ttl time.Duration) error {
	if err := gtr.writeGuard(); err != nil {
		return err
	}
	aKey := fmt.Sprintf("%s:%s", gtr.prefix, key)
	err := gtr.redisClient.Set(ctx, aKey, value, ttl).Err()

	return err
}

// Get 获取单条数据
func (gtr *GTRedis) Get(key string) (string, error) {
	aKey := fmt.Sprintf("%s:%s", gtr.prefix, key)
	val, err := gtr.redisClient.Get(ctx, aKey).Result()

	return val, err
}

// Del 删除单条数据
func (gtr *GTRedis) Del(key string) error {
	if err := gtr.writeGuard(); err != nil {
		return err
	}
	aKey := fmt.Sprintf("%s:%s", gtr.prefix, key)

	err := gtr.redisClient.Del(ctx, aKey).Err()

	return err
}

// Keys 按 pattern 查找键(pattern 支持 redis glob,如 "user:*")。
// 内部走 SCAN 迭代而非 KEYS 命令(KEYS O(N) 阻塞服务端,生产禁用);
// 语义与 KEYS 的差异:非时点快照——迭代期间新增/删除的键可能漏、结果去重后返回。
func (gtr *GTRedis) Keys(key string) ([]string, error) {
	aKey := fmt.Sprintf("%s:%s", gtr.prefix, key)

	var cursor uint64
	seen := make(map[string]struct{})
	var keys []string
	for {
		scanKeys, newCursor, err := gtr.redisClient.Scan(ctx, cursor, aKey, 0).Result()
		if err != nil {
			return nil, err
		}
		for _, k := range scanKeys {
			if _, dup := seen[k]; !dup {
				seen[k] = struct{}{}
				keys = append(keys, k)
			}
		}
		cursor = newCursor
		if cursor == 0 {
			break
		}
	}
	return keys, nil
}
