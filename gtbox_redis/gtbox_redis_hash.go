package gtbox_redis

import (
	"fmt"
	"strings"
)

func (gtr *GTRedis) HGetAll(key string) (map[string]string, error) {
	aKey := fmt.Sprintf("%s:%s", prefix, key)

	val, err := gtr.redisClient.HGetAll(ctx, aKey).Result()
	return val, err
}

// HSet Hash类型-插入单条数据
func (gtr *GTRedis) HSet(key string, subKey string, jsonByte []byte) error {
	aKey := fmt.Sprintf("%s:%s", prefix, key)

	err := gtr.redisClient.HSet(ctx, aKey, subKey, jsonByte).Err()
	return err
}

// HGet Hash类型-获取单条数据
func (gtr *GTRedis) HGet(key string, subKey string) (string, error) {
	aKey := fmt.Sprintf("%s:%s", prefix, key)
	val, a_err := gtr.redisClient.HGet(ctx, aKey, subKey).Result()
	return val, a_err
}

// HDel Hash类型-删除单条数据
func (gtr *GTRedis) HDel(key string, subKey string) error {
	aKey := fmt.Sprintf("%s:%s", prefix, key)

	err := gtr.redisClient.HDel(ctx, aKey, subKey).Err()
	return err
}

// HExists Hash类型-判断是否存在
func (gtr *GTRedis) HExists(key string, subKey string) bool {
	aKey := fmt.Sprintf("%s:%s", prefix, key)
	val, _ := gtr.redisClient.HExists(ctx, aKey, subKey).Result()

	return val
}

func (gtr *GTRedis) HKeys(key string) ([]string, error) {
	aKey := fmt.Sprintf("%s:%s", prefix, key)
	val, err := gtr.redisClient.HKeys(ctx, aKey).Result()

	return val, err
}

// ScanSameLevelKeys 使用 SCAN 命令查找同级的键
func (gtr *GTRedis) ScanSameLevelKeys(key string) ([]string, error) {
	aKey := fmt.Sprintf("%s:%s", prefix, key)

	// 分割 key，并使用最后一个部分的通配符进行模式匹配
	parts := strings.Split(aKey, ":")
	if len(parts) == 0 {
		return nil, fmt.Errorf("invalid key format")
	}
	// 替换最后一个部分为通配符
	parts[len(parts)-1] = "*"
	pattern := strings.Join(parts, ":")
	var cursor uint64
	var keys []string

	// 循环使用 SCAN 命令查找所有匹配的键
	for {
		scanKeys, newCursor, err := gtr.redisClient.Scan(ctx, cursor, pattern, 0).Result()
		if err != nil {
			return nil, err
		}

		keys = append(keys, scanKeys...)
		cursor = newCursor

		if cursor == 0 {
			break
		}
	}

	return keys, nil
}

func (gtr *GTRedis) HScan(key string, cursor uint64, match string, count int64) ([]string, error) {
	aKey := fmt.Sprintf("%s:%s", prefix, key)
	val, _, err := gtr.redisClient.HScan(ctx, aKey, cursor, match, count).Result()
	return val, err
}
