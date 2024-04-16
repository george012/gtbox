package gtbox_redis

import "fmt"

// SAdd 集合--添加数据
func (mbr *MbRedis) SAdd(key1 string, values ...interface{}) error {
	aKey := fmt.Sprintf("%s:%s", prefix, key1)

	err := mbr.redisClient.SAdd(ctx, aKey, values).Err()

	return err
}

// SMembers 集合--获取数据
func (mbr *MbRedis) SMembers(key1 string, values ...interface{}) []string {
	aKey := fmt.Sprintf("%s:%s", prefix, key1)
	val, _ := mbr.redisClient.SMembers(ctx, aKey).Result()

	return val
}

// Scard 集合--获取数据数量
func (mbr *MbRedis) Scard(key1 string, key2 string) (Cnt int64) {
	aKey := fmt.Sprintf("%s:%s", prefix, key1)
	val, _ := mbr.redisClient.SCard(ctx, aKey).Result()

	return val
}
