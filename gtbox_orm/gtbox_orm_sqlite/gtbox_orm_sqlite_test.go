package gtbox_orm_sqlite

import (
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

type ormTestUser struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

func newTestSqlite(t *testing.T, fileName string) *GTORMSqlite {
	t.Helper()

	aSqlite, err := New(GTORMSqliteConfig{DBFilePath: filepath.Join(t.TempDir(), fileName)})
	if err != nil {
		t.Fatalf("New() err=%v", err)
	}
	t.Cleanup(func() { _ = aSqlite.Close() })

	if err := aSqlite.SqliteDB.AutoMigrate(&ormTestUser{}); err != nil {
		t.Fatalf("AutoMigrate() err=%v", err)
	}
	return aSqlite
}

// TestGTORMSqliteConfigValidate 库文件路径缺失、池参数非法必须显式拒绝
func TestGTORMSqliteConfigValidate(t *testing.T) {
	cases := []struct {
		name    string
		cfg     GTORMSqliteConfig
		wantErr bool
	}{
		{"合法配置", GTORMSqliteConfig{DBFilePath: "app.db"}, false},
		{"缺 DBFilePath", GTORMSqliteConfig{}, true},
		{"ConnMaxOpen 负数", GTORMSqliteConfig{DBFilePath: "app.db", ConnMaxOpen: -1}, true},
		{"ConnMaxIdle 超过 ConnMaxOpen", GTORMSqliteConfig{DBFilePath: "app.db", ConnMaxOpen: 1, ConnMaxIdle: 4}, true},
		{"ConnMaxIdleTime 负数", GTORMSqliteConfig{DBFilePath: "app.db", ConnMaxIdleTime: -time.Second}, true},
	}

	for _, aCase := range cases {
		t.Run(aCase.name, func(t *testing.T) {
			err := aCase.cfg.validate()
			if (err != nil) != aCase.wantErr {
				t.Fatalf("validate() err=%v, wantErr=%v", err, aCase.wantErr)
			}
			if aCase.wantErr {
				if _, newErr := New(aCase.cfg); newErr == nil {
					t.Fatalf("New() 应在参数校验阶段失败")
				}
			}
		})
	}
}

// TestGTORMSqliteConfigPoolParams 池参数零值取 sqlite 单写默认(1/1),不套用 mysql 那组
func TestGTORMSqliteConfigPoolParams(t *testing.T) {
	maxOpen, maxIdle, idleTime := (&GTORMSqliteConfig{}).poolParams()
	if maxOpen != defaultSqliteConnMaxOpen || maxIdle != defaultSqliteConnMaxIdle || idleTime != 0 {
		t.Fatalf("零值池参数=%d/%d/%s, want %d/%d/0", maxOpen, maxIdle, idleTime, defaultSqliteConnMaxOpen, defaultSqliteConnMaxIdle)
	}
}

// TestGTORMSqliteMultiInstanceIsolation 两个 New 出来的实例各自独立库文件、独立连接池,数据互不串
func TestGTORMSqliteMultiInstanceIsolation(t *testing.T) {
	orderDB := newTestSqlite(t, "order.db")
	auditDB := newTestSqlite(t, "audit.db")

	if err := orderDB.InsertData(&ormTestUser{Name: "order-user"}); err != nil {
		t.Fatalf("orderDB.InsertData() err=%v", err)
	}

	var auditRows []ormTestUser
	if err := auditDB.SqliteDB.Find(&auditRows).Error; err != nil {
		t.Fatalf("auditDB.Find() err=%v", err)
	}
	if len(auditRows) != 0 {
		t.Fatalf("实例间数据串了,auditDB 行数=%d", len(auditRows))
	}

	var orderRows []ormTestUser
	if err := orderDB.SqliteDB.Find(&orderRows).Error; err != nil {
		t.Fatalf("orderDB.Find() err=%v", err)
	}
	if len(orderRows) != 1 || orderRows[0].Name != "order-user" {
		t.Fatalf("orderDB 数据不对: %+v", orderRows)
	}
}

// TestGTORMSqliteConcurrentInsert 去掉包内互斥锁后,并发写靠连接池串行,不应出现 database is locked
func TestGTORMSqliteConcurrentInsert(t *testing.T) {
	aSqlite := newTestSqlite(t, "concurrent.db")

	const workerCount = 32
	var wg sync.WaitGroup
	errCh := make(chan error, workerCount)
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			if err := aSqlite.InsertData(&ormTestUser{Name: "user-" + string(rune('a'+idx%26)) + string(rune('0'+idx/26))}); err != nil {
				errCh <- err
			}
		}(i)
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		t.Fatalf("并发 InsertData err=%v", err)
	}

	var rows []ormTestUser
	if err := aSqlite.SqliteDB.Find(&rows).Error; err != nil {
		t.Fatalf("Find() err=%v", err)
	}
	if len(rows) != workerCount {
		t.Fatalf("并发写入行数=%d, want %d", len(rows), workerCount)
	}
}

// TestGTORMSqliteCloseAndReopenGuard Close 后方法显式报未打开;已打开的句柄拒绝重复 Open
func TestGTORMSqliteCloseAndReopenGuard(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "guard.db")

	aSqlite, err := New(GTORMSqliteConfig{DBFilePath: dbPath})
	if err != nil {
		t.Fatalf("New() err=%v", err)
	}

	if err := aSqlite.OpenWithConfig(GTORMSqliteConfig{DBFilePath: dbPath}); err == nil || !strings.Contains(err.Error(), "already opened") {
		t.Fatalf("重复 Open err=%v, want already opened", err)
	}

	if err := aSqlite.Close(); err != nil {
		t.Fatalf("Close() err=%v", err)
	}
	if err := aSqlite.InsertData(&ormTestUser{Name: "x"}); err == nil || !strings.Contains(err.Error(), "not opened") {
		t.Fatalf("Close 后 InsertData err=%v, want not opened", err)
	}
	if err := aSqlite.Close(); err != nil {
		t.Fatalf("重复 Close() err=%v", err)
	}

	// Close 之后可以重新打开,句柄复用
	if err := aSqlite.OpenWithConfig(GTORMSqliteConfig{DBFilePath: dbPath}); err != nil {
		t.Fatalf("Close 后重新 Open err=%v", err)
	}
	_ = aSqlite.Close()
}

// TestGTORMSqliteInstanceDefaultSingleton 默认单例原有入口 OpenSqlite 行为不变
func TestGTORMSqliteInstanceDefaultSingleton(t *testing.T) {
	if Instance() != Instance() {
		t.Fatalf("Instance() 返回了不同句柄")
	}

	Instance().OpenSqlite(filepath.Join(t.TempDir(), "singleton.db"))
	if Instance().SqliteError != nil {
		t.Fatalf("OpenSqlite() err=%v", Instance().SqliteError)
	}
	t.Cleanup(func() { _ = Instance().Close() })

	if err := Instance().SqliteDB.AutoMigrate(&ormTestUser{}); err != nil {
		t.Fatalf("AutoMigrate() err=%v", err)
	}
	if err := Instance().InsertData(&ormTestUser{Name: "singleton-user"}); err != nil {
		t.Fatalf("InsertData() err=%v", err)
	}

	queried := ormTestUser{Name: "singleton-user"}
	if err := Instance().QueryData(&queried); err != nil {
		t.Fatalf("QueryData() err=%v", err)
	}
	if queried.ID == 0 {
		t.Fatalf("QueryData 没查回数据: %+v", queried)
	}
}
