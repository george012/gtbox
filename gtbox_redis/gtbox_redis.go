package gtbox_redis

import (
	"context"
	"fmt"
	"github.com/george012/gtbox/gtbox_log"
	"github.com/redis/go-redis/v9"
)

var (
	isSetup    bool
	OwnerRedis *MbRedis
	ctx        = context.Background()
	prefix     string
)

type RedisConfig struct {
	Addr       string `yaml:"addr" json:"addr"`              // address
	Pwd        string `yaml:"pwd" json:"pwd"`                // pwd
	SocketBuck int    `yaml:"socketBuck" json:"socket_buck"` // 插槽
}

type MbRedis struct {
	cfg         *RedisConfig
	redisClient *redis.Client
}

// SetupRedisConnection 初始化Redis连接
func SetupRedisConnection(redisCfg RedisConfig, prefixStr string) (success bool) {
	if isSetup == false {
		OwnerRedis = &MbRedis{
			cfg: &redisCfg,
			redisClient: redis.NewClient(&redis.Options{
				Addr:     redisCfg.Addr,
				Password: redisCfg.Pwd,        // no password set
				DB:       redisCfg.SocketBuck, // use default DB
			}),
		}

		prefix = prefixStr
		err := OwnerRedis.redisClient.Set(ctx, "hello", "helloValue", 0).Err()
		if err != nil {
			gtbox_log.LogErrorf("[redis setup] error [%s]", err)
			isSetup = false
		} else {
			gtbox_log.LogInfof("[redis setup] [%s]", "success")
			isSetup = true
		}
	}
	return isSetup
}

// Set 插入单条数据
func (mbr *MbRedis) Set(key string, value string) error {
	aKey := fmt.Sprintf("%s:%s", prefix, key)
	err := mbr.redisClient.Set(ctx, aKey, value, 0).Err()

	return err
}

// Get 获取单条数据
func (mbr *MbRedis) Get(key string) (string, error) {
	aKey := fmt.Sprintf("%s:%s", prefix, key)
	val, err := mbr.redisClient.Get(ctx, aKey).Result()

	return val, err
}

// Del 删除单条数据
func (mbr *MbRedis) Del(key string) error {
	aKey := fmt.Sprintf("%s:%s", prefix, key)

	err := mbr.redisClient.Del(ctx, aKey).Err()

	return err
}

// Keys 删除单条数据
func (mbr *MbRedis) Keys(key string) ([]string, error) {
	aKey := fmt.Sprintf("%s:%s", prefix, key)
	val, err := mbr.redisClient.Keys(ctx, aKey).Result()

	return val, err
}
