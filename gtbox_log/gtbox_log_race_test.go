package gtbox_log

import (
	"sync"
	"testing"
)

// TestSetupDefaultLog_ConcurrentLazyInit_NoRace 坐实包级 LogXxxf 入口的懒加载并发安全。
//
// 修复前 setupDefaultLog 是无锁的 check-then-act（if setupComplete==false && mainLog==nil），
// 多 goroutine 首次并发调日志会同时读写包级 mainLog → data race（go test -race 可复现）。
// 修复后用 mainLogOnce 收口，本测试在 -race 下应干净通过、且 mainLog 被创建一次。
//
// 说明：本包此前零测试，故 mainLogOnce 在本进程内尚未触发，这里正是对「并发首次初始化」
// 这一真实竞争窗口的覆盖。enableSaveLogFile 默认 false → 不落文件、不建目录，只走 stdout。
func TestSetupDefaultLog_ConcurrentLazyInit_NoRace(t *testing.T) {
	// 并发前单线程把配置写好（不触发 mainLog 创建），避免落到默认 "/" 目录。
	instanceConfig().productName = "gtbox_log_racetest"
	instanceConfig().productLogDir = t.TempDir()

	const n = 32
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			LogInfof("concurrent lazy-init log %d", i)
		}(i)
	}
	wg.Wait()

	if mainLog == nil {
		t.Fatal("并发懒加载后 mainLog 仍为 nil")
	}
}
