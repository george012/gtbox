package gtbox_redis

import "fmt"

// SAdd 集合--添加数据
func (gtr *GTRedis) SAdd(key1 string, values ...interface{}) error {
	if err := gtr.writeGuard(); err != nil {
		return err
	}
	aKey := fmt.Sprintf("%s:%s", gtr.prefix, key1)

	err := gtr.redisClient.SAdd(ctx, aKey, values...).Err()

	return err
}

// SMembers 集合--获取数据
func (gtr *GTRedis) SMembers(key1 string, values ...interface{}) ([]string, error) {
	aKey := fmt.Sprintf("%s:%s", gtr.prefix, key1)
	val, err := gtr.redisClient.SMembers(ctx, aKey).Result()

	return val, err
}

// Scard 集合--获取数据数量
func (gtr *GTRedis) Scard(key1 string, key2 string) (Cnt int64, err error) {
	aKey := fmt.Sprintf("%s:%s", gtr.prefix, key1)
	val, err := gtr.redisClient.SCard(ctx, aKey).Result()

	return val, err
}
