package gtbox_redis

import (
	"runtime"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
)

// newTestRedis 起一个 in-process miniredis 并建 GTRedis 实例
func newTestRedis(t *testing.T, cfgMutate func(*RedisConfig), prefix string) (*miniredis.Miniredis, *GTRedis) {
	t.Helper()
	mr := miniredis.RunT(t)
	cfg := RedisConfig{Addr: mr.Addr()}
	if cfgMutate != nil {
		cfgMutate(&cfg)
	}
	gtr, err := New(cfg, prefix)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = gtr.Close() })
	return mr, gtr
}

func TestNewParamValidation(t *testing.T) {
	if _, err := New(RedisConfig{}, "p"); err == nil {
		t.Fatal("empty Addr should be rejected")
	}
	if _, err := New(RedisConfig{Addr: "127.0.0.1:0", PoolSizeMultiplier: -1}, "p"); err == nil {
		t.Fatal("negative PoolSizeMultiplier should be rejected")
	}
}

func TestPoolSizeMultiplierPassthrough(t *testing.T) {
	_, gtr := newTestRedis(t, func(c *RedisConfig) { c.PoolSizeMultiplier = 3 }, "p")
	want := 3 * 10 * runtime.GOMAXPROCS(0)
	if got := gtr.redisClient.Options().PoolSize; got != want {
		t.Fatalf("PoolSize got %d, want %d", got, want)
	}
	// 0 = go-redis 默认(不显式设置)
	_, gtrDefault := newTestRedis(t, nil, "p")
	if got := gtrDefault.redisClient.Options().PoolSize; got != 10*runtime.GOMAXPROCS(0) {
		t.Fatalf("default PoolSize got %d, want %d", got, 10*runtime.GOMAXPROCS(0))
	}
}

// TestSetupRedisConnectionConcurrent 并发 Setup 收口验证(go test -race)
func TestSetupRedisConnectionConcurrent(t *testing.T) {
	mr := miniredis.RunT(t)
	cfg := RedisConfig{Addr: mr.Addr()}

	var wg sync.WaitGroup
	results := make([]bool, 16)
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx] = SetupRedisConnection(cfg, "race_test")
		}(i)
	}
	wg.Wait()

	if OwnerRedis == nil {
		t.Fatal("OwnerRedis should be initialized")
	}
	for i, ok := range results {
		if !ok {
			t.Fatalf("goroutine %d: Setup should report success", i)
		}
	}
	if err := OwnerRedis.Set("k", "v"); err != nil {
		t.Fatalf("singleton Set failed: %v", err)
	}
	if v, err := OwnerRedis.Get("k"); err != nil || v != "v" {
		t.Fatalf("singleton Get got (%q, %v), want (\"v\", nil)", v, err)
	}
}

func TestMultiInstanceIsolation(t *testing.T) {
	_, r1 := newTestRedis(t, nil, "app1")
	_, r2 := newTestRedis(t, nil, "app2")

	if err := r1.Set("k", "v1"); err != nil {
		t.Fatal(err)
	}
	if err := r2.Set("k", "v2"); err != nil {
		t.Fatal(err)
	}
	if v, _ := r1.Get("k"); v != "v1" {
		t.Fatalf("r1 got %q, want v1", v)
	}
	if v, _ := r2.Get("k"); v != "v2" {
		t.Fatalf("r2 got %q, want v2", v)
	}
}

func TestReadOnlyGuard(t *testing.T) {
	_, ro := newTestRedis(t, func(c *RedisConfig) { c.ReadOnly = true }, "ro")

	if err := ro.Set("k", "v"); err == nil {
		t.Fatal("Set on read-only instance should be rejected")
	}
	if err := ro.SetEX("k", "v", time.Minute); err == nil {
		t.Fatal("SetEX on read-only instance should be rejected")
	}
	if err := ro.Del("k"); err == nil {
		t.Fatal("Del on read-only instance should be rejected")
	}
	if err := ro.HSet("h", "f", []byte("v")); err == nil {
		t.Fatal("HSet on read-only instance should be rejected")
	}
	if err := ro.HDel("h", "f"); err == nil {
		t.Fatal("HDel on read-only instance should be rejected")
	}
	if err := ro.SAdd("s", "m"); err == nil {
		t.Fatal("SAdd on read-only instance should be rejected")
	}
	if _, err := ro.SetNX("k", "v", time.Minute); err == nil {
		t.Fatal("SetNX on read-only instance should be rejected")
	}
	if _, err := ro.HSetNX("h", "f", []byte("v")); err == nil {
		t.Fatal("HSetNX on read-only instance should be rejected")
	}
	// 读路径正常(键不存在返回 redis.Nil 属正常读语义)
	if _, err := ro.Keys("*"); err != nil {
		t.Fatalf("read on read-only instance should work: %v", err)
	}
}

func TestSetEXTTL(t *testing.T) {
	mr, gtr := newTestRedis(t, nil, "ttl")
	if err := gtr.SetEX("k", "v", time.Minute); err != nil {
		t.Fatal(err)
	}
	if ttl := mr.TTL("ttl:k"); ttl != time.Minute {
		t.Fatalf("TTL got %v, want 1m", ttl)
	}
	// Set 无过期
	if err := gtr.Set("k2", "v"); err != nil {
		t.Fatal(err)
	}
	if ttl := mr.TTL("ttl:k2"); ttl != 0 {
		t.Fatalf("persistent key TTL got %v, want 0", ttl)
	}
}

func TestSetNX(t *testing.T) {
	mr, gtr := newTestRedis(t, nil, "nx")

	ok, err := gtr.SetNX("lock", "owner-1", time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("first SetNX should write")
	}
	if ttl := mr.TTL("nx:lock"); ttl != time.Minute {
		t.Fatalf("TTL got %v, want 1m", ttl)
	}

	// 键已存在:返回 false 且不覆盖原值,不算错误
	ok, err = gtr.SetNX("lock", "owner-2", time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("second SetNX should not write")
	}
	val, err := gtr.Get("lock")
	if err != nil {
		t.Fatal(err)
	}
	if val != "owner-1" {
		t.Fatalf("value got %q, want owner-1", val)
	}

	// ttl 为 0 = 持久键
	if ok, err = gtr.SetNX("persistent", "v", 0); err != nil || !ok {
		t.Fatalf("SetNX persistent ok=%v err=%v", ok, err)
	}
	if ttl := mr.TTL("nx:persistent"); ttl != 0 {
		t.Fatalf("persistent key TTL got %v, want 0", ttl)
	}

	// 过期后同一把锁可以重新拿到
	mr.FastForward(time.Minute)
	if ok, err = gtr.SetNX("lock", "owner-2", time.Minute); err != nil || !ok {
		t.Fatalf("SetNX after expiry ok=%v err=%v", ok, err)
	}
}

func TestHSetNX(t *testing.T) {
	_, gtr := newTestRedis(t, nil, "hnx")

	ok, err := gtr.HSetNX("h", "f", []byte("v1"))
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("first HSetNX should write")
	}

	// 字段已存在:返回 false 且不覆盖,不算错误
	ok, err = gtr.HSetNX("h", "f", []byte("v2"))
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("second HSetNX should not write")
	}
	val, err := gtr.HGet("h", "f")
	if err != nil {
		t.Fatal(err)
	}
	if val != "v1" {
		t.Fatalf("field value got %q, want v1", val)
	}

	// 同键不同字段互不影响
	if ok, err = gtr.HSetNX("h", "f2", []byte("v3")); err != nil || !ok {
		t.Fatalf("HSetNX other field ok=%v err=%v", ok, err)
	}
}

func TestKeysScanGlob(t *testing.T) {
	_, gtr := newTestRedis(t, nil, "app")
	for _, k := range []string{"user:1", "user:2", "order:1"} {
		if err := gtr.Set(k, "v"); err != nil {
			t.Fatal(err)
		}
	}
	keys, err := gtr.Keys("user:*")
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(keys)
	want := []string{"app:user:1", "app:user:2"}
	if len(keys) != len(want) || keys[0] != want[0] || keys[1] != want[1] {
		t.Fatalf("Keys got %v, want %v", keys, want)
	}
}

func TestScanSameLevelKeys(t *testing.T) {
	_, gtr := newTestRedis(t, nil, "app")
	for _, k := range []string{"group:a", "group:b", "group:c"} {
		if err := gtr.Set(k, "v"); err != nil {
			t.Fatal(err)
		}
	}
	keys, err := gtr.ScanSameLevelKeys("group:a")
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 3 {
		t.Fatalf("same-level got %d keys %v, want 3", len(keys), keys)
	}
	seen := make(map[string]struct{})
	for _, k := range keys {
		if _, dup := seen[k]; dup {
			t.Fatalf("duplicate key %q in result", k)
		}
		seen[k] = struct{}{}
	}
}

func TestHashOps(t *testing.T) {
	_, gtr := newTestRedis(t, nil, "h")
	if err := gtr.HSet("obj", "f1", []byte("v1")); err != nil {
		t.Fatal(err)
	}
	if v, err := gtr.HGet("obj", "f1"); err != nil || v != "v1" {
		t.Fatalf("HGet got (%q, %v)", v, err)
	}
	exists, err := gtr.HExists("obj", "f1")
	if err != nil || !exists {
		t.Fatalf("HExists got (%v, %v), want (true, nil)", exists, err)
	}
	exists, err = gtr.HExists("obj", "missing")
	if err != nil || exists {
		t.Fatalf("HExists missing got (%v, %v), want (false, nil)", exists, err)
	}
	if n, err := gtr.HLen("obj"); err != nil || n != 1 {
		t.Fatalf("HLen got (%d, %v)", n, err)
	}
}

func TestHScanCursor(t *testing.T) {
	_, gtr := newTestRedis(t, nil, "h")
	for i := 0; i < 30; i++ {
		if err := gtr.HSet("big", string(rune('a'+i%26))+string(rune('0'+i/26)), []byte("v")); err != nil {
			t.Fatal(err)
		}
	}
	// 按游标翻页收齐全部 field(HScan 返回 field/value 交替对)
	fields := make(map[string]struct{})
	var cursor uint64
	for {
		pairs, next, err := gtr.HScan("big", cursor, "*", 10)
		if err != nil {
			t.Fatal(err)
		}
		for i := 0; i < len(pairs); i += 2 {
			fields[pairs[i]] = struct{}{}
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	if len(fields) != 30 {
		t.Fatalf("paged HScan collected %d fields, want 30", len(fields))
	}
}

func TestSetOps(t *testing.T) {
	_, gtr := newTestRedis(t, nil, "s")
	if err := gtr.SAdd("tags", "a", "b", "c"); err != nil {
		t.Fatal(err)
	}
	members, err := gtr.SMembers("tags")
	if err != nil {
		t.Fatal(err)
	}
	if len(members) != 3 {
		t.Fatalf("SMembers got %v, want 3 members", members)
	}
	cnt, err := gtr.Scard("tags", "")
	if err != nil || cnt != 3 {
		t.Fatalf("Scard got (%d, %v), want (3, nil)", cnt, err)
	}
}

// TestErrSurfacedAfterClose 吞错修复验证:连接关闭后错误必须浮出而非零值
func TestErrSurfacedAfterClose(t *testing.T) {
	_, gtr := newTestRedis(t, nil, "e")
	if err := gtr.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := gtr.HExists("k", "f"); err == nil {
		t.Fatal("HExists after Close should return error")
	}
	if _, err := gtr.SMembers("k"); err == nil {
		t.Fatal("SMembers after Close should return error")
	}
	if _, err := gtr.Scard("k", ""); err == nil {
		t.Fatal("Scard after Close should return error")
	}
}
