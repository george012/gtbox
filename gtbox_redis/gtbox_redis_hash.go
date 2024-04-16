package gtbox_redis

import (
	"fmt"
)

func (mbr *MbRedis) HGetAll(key string) (map[string]string, error) {
	aKey := fmt.Sprintf("%s:%s", prefix, key)

	val, err := mbr.redisClient.HGetAll(ctx, aKey).Result()
	return val, err
}

// HSet Hash类型-插入单条数据
func (mbr *MbRedis) HSet(key string, subKey string, jsonByte []byte) error {
	aKey := fmt.Sprintf("%s:%s", prefix, key)

	err := mbr.redisClient.HSet(ctx, aKey, subKey, jsonByte).Err()
	return err
}

// HGet Hash类型-获取单条数据
func (mbr *MbRedis) HGet(key string, subKey string) (string, error) {
	aKey := fmt.Sprintf("%s:%s", prefix, key)
	val, a_err := mbr.redisClient.HGet(ctx, aKey, subKey).Result()
	return val, a_err
}

// HDel Hash类型-删除单条数据
func (mbr *MbRedis) HDel(key string, subKey string) error {
	aKey := fmt.Sprintf("%s:%s", prefix, key)

	err := mbr.redisClient.HDel(ctx, aKey, subKey).Err()
	return err
}

// HExists Hash类型-判断是否存在
func (mbr *MbRedis) HExists(key string, subKey string) bool {
	aKey := fmt.Sprintf("%s:%s", prefix, key)
	val, _ := mbr.redisClient.HExists(ctx, aKey, subKey).Result()

	return val
}

func (mbr *MbRedis) HKeys(key string) ([]string, error) {
	aKey := fmt.Sprintf("%s:%s", prefix, key)
	val, err := mbr.redisClient.HKeys(ctx, aKey).Result()

	return val, err
}
