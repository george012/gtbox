/*
Package gtbox_orm_sqlite en: sqlite handle, zh-cn: sqlite 封装处理

默认单例 + 多实例并存(对齐 gtbox_orm_mysql / gtbox_redis 模式):
  - 单例:Instance() 取默认句柄,OpenSqlite / OpenWithConfig 打开连接
  - 多实例:New 创建独立句柄(独立 gorm.DB / 独立连接池),方法集与单例完全一致,
    用于一个进程同时接多个 sqlite 库文件
*/
package gtbox_orm_sqlite

import (
	"fmt"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GTORMSqlite struct {
	SqliteDB    *gorm.DB
	SqliteError error
	mux         sync.RWMutex
}

var (
	gtSqliteOnce   sync.Once
	sqliteInstance *GTORMSqlite
)

// Instance 默认单例句柄
func Instance() *GTORMSqlite {
	gtSqliteOnce.Do(func() {
		sqliteInstance = &GTORMSqlite{}
	})
	return sqliteInstance
}

// New 创建独立 sqlite 实例:独立 gorm.DB、独立连接池,与 Instance() 默认单例互不影响。
// error 既包含参数校验失败,也包含首次打开库文件失败(gorm.Open 内部默认带 Ping)。
func New(cfg GTORMSqliteConfig) (*GTORMSqlite, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	db, err := openSqliteDB(&cfg)
	if err != nil {
		return nil, err
	}
	return &GTORMSqlite{SqliteDB: db}, nil
}

// openSqliteDB 打开库文件 + 连接池配置,New 与 OpenSqlite / OpenWithConfig 共用这一条路径
func openSqliteDB(cfg *GTORMSqliteConfig) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.DBFilePath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	maxOpen, maxIdle, idleTime := cfg.poolParams()
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)
	if idleTime > 0 {
		sqlDB.SetConnMaxIdleTime(idleTime)
	}
	return db, nil
}

// OpenWithConfig 用配置打开当前句柄的连接(默认单例走 Instance().OpenWithConfig(cfg))。
// 句柄已打开时直接返回 error:替换句柄会让已持有 SqliteDB 的调用方指向旧池,旧池又无人关闭。
func (gtSqlite *GTORMSqlite) OpenWithConfig(cfg GTORMSqliteConfig) error {
	gtSqlite.mux.Lock()
	defer gtSqlite.mux.Unlock()

	if err := cfg.validate(); err != nil {
		gtSqlite.SqliteError = err
		return err
	}

	if gtSqlite.SqliteDB != nil {
		err := fmt.Errorf("gtbox_orm_sqlite: instance already opened, close it before reopen")
		gtSqlite.SqliteError = err
		return err
	}

	db, err := openSqliteDB(&cfg)
	gtSqlite.SqliteDB, gtSqlite.SqliteError = db, err
	return err
}

// OpenSqlite 打开当前句柄的库文件(原有入口,签名不变)。
// 连接池取包内默认值(ConnMaxOpen 1 / ConnMaxIdle 1);要自定义池参数用 OpenWithConfig。
func (gtSqlite *GTORMSqlite) OpenSqlite(sqlitePath string) {
	gtSqlite.mux.Lock()
	defer gtSqlite.mux.Unlock()

	cfg := &GTORMSqliteConfig{DBFilePath: sqlitePath}
	gtSqlite.SqliteDB, gtSqlite.SqliteError = openSqliteDB(cfg)
	if gtSqlite.SqliteError != nil {
		println("连接数据库失败==", gtSqlite.SqliteDB, gtSqlite.SqliteError)
	}
}

// Close 关闭当前句柄的连接池。默认单例一般随进程生命周期存在,主要给 New 出来的实例用。
func (gtSqlite *GTORMSqlite) Close() error {
	gtSqlite.mux.Lock()
	defer gtSqlite.mux.Unlock()

	if gtSqlite.SqliteDB == nil {
		return nil
	}
	sqlDB, err := gtSqlite.SqliteDB.DB()
	if err != nil {
		return err
	}
	if err := sqlDB.Close(); err != nil {
		return err
	}
	gtSqlite.SqliteDB = nil
	return nil
}

// sqliteDB 取当前句柄的 gorm.DB。只用读锁护住指针本身(与 OpenXxx / Close 的写互斥),
// 取到后即释放锁再跑 SQL:写的串行化交给连接池(默认 MaxOpen 1),不再靠包内互斥锁——
// 包内锁只护得住本包这两个方法,护不住调用方直接拿 SqliteDB 发出的 SQL,连接池两边都护得住。
func (gtSqlite *GTORMSqlite) sqliteDB() (*gorm.DB, error) {
	gtSqlite.mux.RLock()
	defer gtSqlite.mux.RUnlock()

	if gtSqlite.SqliteDB == nil {
		return nil, fmt.Errorf("gtbox_orm_sqlite: instance not opened")
	}
	return gtSqlite.SqliteDB, nil
}

// InsertData 数据不存在时插入
func (gtSqlite *GTORMSqlite) InsertData(dataModel interface{}) error {
	db, err := gtSqlite.sqliteDB()
	if err != nil {
		return err
	}

	result := db.Where(dataModel).Limit(1).Find(dataModel)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		if cuerr := db.Create(dataModel).Error; cuerr != nil {
			return cuerr
		}
	}
	return nil
}

// QueryData 按 dataModel 条件查询
func (gtSqlite *GTORMSqlite) QueryData(dataModel interface{}, conditions ...interface{}) error {
	db, err := gtSqlite.sqliteDB()
	if err != nil {
		return err
	}
	return db.Where(dataModel, conditions...).Find(dataModel).Error
}
