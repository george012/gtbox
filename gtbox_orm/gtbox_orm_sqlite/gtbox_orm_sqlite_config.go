package gtbox_orm_sqlite

import (
	"fmt"
	"time"
)

// 连接池默认值。sqlite 的写锁在文件级,多开连接并不会让写并行:
// 多连接并发写是在 driver 的 busy_timeout 里排队等锁(mattn/go-sqlite3 默认 5000ms,
// sqlite3.go:1201),等超了才报 database is locked。默认单连接把排队挪到连接池,
// 等待可预期、不会超时报错。要靠多连接吃并发读的,自己调大 ConnMaxOpen,
// 并在 DBFilePath 上带 _journal_mode=WAL。
const (
	defaultSqliteConnMaxOpen = 1
	defaultSqliteConnMaxIdle = 1
)

// GTORMSqliteConfig 一个 sqlite 数据源的文件路径与连接池参数。
// DBFilePath 无默认值,零值直接报错;其内容原样交给 driver,可带 driver 的 DSN 参数,
// 例如 "app.db?_journal_mode=WAL&_busy_timeout=5000" 或 "file::memory:?cache=shared"。
// 池参数零值 = 取上面那组默认值;ConnMaxIdleTime 零值 = 不回收空闲连接(sqlite 连接就是个文件句柄,常驻更省事)。
type GTORMSqliteConfig struct {
	DBFilePath      string        `yaml:"db_file_path" json:"db_file_path"`
	ConnMaxOpen     int           `yaml:"conn_max_open" json:"conn_max_open"`
	ConnMaxIdle     int           `yaml:"conn_max_idle" json:"conn_max_idle"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time" json:"conn_max_idle_time"`
}

// validate 参数校验;ConnMaxIdle > ConnMaxOpen 时 database/sql 会静默把 idle 压到 open,
// 这里改为直接拒绝,避免配置写错却查不出来。
func (cfg *GTORMSqliteConfig) validate() error {
	if cfg.DBFilePath == "" {
		return fmt.Errorf("gtbox_orm_sqlite: DBFilePath is required")
	}
	if cfg.ConnMaxOpen < 0 {
		return fmt.Errorf("gtbox_orm_sqlite: ConnMaxOpen must be >= 0, got %d", cfg.ConnMaxOpen)
	}
	if cfg.ConnMaxIdle < 0 {
		return fmt.Errorf("gtbox_orm_sqlite: ConnMaxIdle must be >= 0, got %d", cfg.ConnMaxIdle)
	}
	if cfg.ConnMaxIdleTime < 0 {
		return fmt.Errorf("gtbox_orm_sqlite: ConnMaxIdleTime must be >= 0, got %s", cfg.ConnMaxIdleTime)
	}
	if cfg.ConnMaxOpen > 0 && cfg.ConnMaxIdle > cfg.ConnMaxOpen {
		return fmt.Errorf("gtbox_orm_sqlite: ConnMaxIdle(%d) must be <= ConnMaxOpen(%d)", cfg.ConnMaxIdle, cfg.ConnMaxOpen)
	}
	return nil
}

// poolParams 返回补齐默认值后的池参数
func (cfg *GTORMSqliteConfig) poolParams() (maxOpen int, maxIdle int, idleTime time.Duration) {
	maxOpen, maxIdle, idleTime = cfg.ConnMaxOpen, cfg.ConnMaxIdle, cfg.ConnMaxIdleTime
	if maxOpen == 0 {
		maxOpen = defaultSqliteConnMaxOpen
	}
	if maxIdle == 0 {
		maxIdle = defaultSqliteConnMaxIdle
	}
	return maxOpen, maxIdle, idleTime
}
