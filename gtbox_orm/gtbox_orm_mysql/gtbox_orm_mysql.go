/*
Package gtbox_orm_mysql en: mysql handle, zh-cn: mysql 封装处理

默认单例 + 多实例并存(对齐 gtbox_redis 模式):
  - 单例:Instance() 取默认句柄,OPenMysql / OpenWithConfig 打开连接,业务侧直接复用
  - 多实例:New 创建独立句柄(独立 gorm.DB / 独立连接池),方法集与单例完全一致,
    用于一个进程接多个 mysql 数据源
*/
package gtbox_orm_mysql

import (
	"errors"
	"fmt"
	"sync"

	"github.com/george012/gtbox/gtbox_orm/gtbox_orm_config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GTORMMysql struct {
	MysqlDB    *gorm.DB
	MysqlError error
	mux        sync.RWMutex
}

var (
	mysqlOnce     sync.Once
	mysqlInstance *GTORMMysql
)

// Instance 默认单例句柄。纯取句柄,不动 MysqlError——原先每次调用都在锁外把 MysqlError 置 nil,
// 与打开路径的写构成 data race;错误状态归打开动作所有,由 OPenMysql / OpenWithConfig 写。
func Instance() *GTORMMysql {
	mysqlOnce.Do(func() {
		mysqlInstance = &GTORMMysql{}
	})
	return mysqlInstance
}

// New 创建独立 mysql 实例:独立 gorm.DB、独立连接池,与 Instance() 默认单例互不影响。
// 与 gtbox_redis.New 的差异:gorm.Open 内部默认带 Ping,故本函数的 error 既包含参数校验失败,
// 也包含首次连接失败——连不上即刻暴露,不留到第一条 SQL。
func New(cfg GTORMMysqlConfig) (*GTORMMysql, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	db, err := openMysqlDB(&cfg)
	if err != nil {
		return nil, err
	}
	return &GTORMMysql{MysqlDB: db}, nil
}

// openMysqlDB 建连 + 连接池配置,New 与 OPenMysql / OpenWithConfig 共用这一条路径
func openMysqlDB(cfg *GTORMMysqlConfig) (*gorm.DB, error) {
	loc, err := mysqlTimeZoneLocation(cfg.TimeZone)
	if err != nil {
		return nil, fmt.Errorf("gtbox_orm_mysql: load timezone failed, %w", err)
	}

	db, err := gorm.Open(mysql.Open(buildMysqlDSN(cfg, loc)), &gorm.Config{
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
	sqlDB.SetConnMaxIdleTime(idleTime)
	if cfg.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}
	return db, nil
}

// OpenWithConfig 用配置打开当前句柄的连接(默认单例走 Instance().OpenWithConfig(cfg))。
// 句柄已打开时直接返回 error:替换句柄会让已持有 MysqlDB 的调用方指向旧池,旧池又无人关闭。
func (aMysql *GTORMMysql) OpenWithConfig(cfg GTORMMysqlConfig) error {
	aMysql.mux.Lock()
	defer aMysql.mux.Unlock()

	if err := cfg.validate(); err != nil {
		aMysql.MysqlError = err
		return err
	}

	if aMysql.MysqlDB != nil {
		err := fmt.Errorf("gtbox_orm_mysql: instance already opened, close it before reopen")
		aMysql.MysqlError = err
		return err
	}

	db, err := openMysqlDB(&cfg)
	aMysql.MysqlDB, aMysql.MysqlError = db, err
	return err
}

// OPenMysql 打开当前句柄的连接(原有入口,签名与回调契约不变)。
// 连接池取包内默认值(MaxOpen 5 / MaxIdle 2 / IdleTime 1min);要自定义池参数用 OpenWithConfig。
func (aMysql *GTORMMysql) OPenMysql(dbUser string, dbPwd string, dbName string, dbAddress string, dbPort int, timeZone gtbox_orm_config.GTORMTimeZone, endFunc func(err error)) {
	aMysql.mux.Lock()
	defer aMysql.mux.Unlock()

	cfg := &GTORMMysqlConfig{
		User:     dbUser,
		Pwd:      dbPwd,
		DBName:   dbName,
		Host:     dbAddress,
		Port:     dbPort,
		TimeZone: timeZone,
	}
	aMysql.MysqlDB, aMysql.MysqlError = openMysqlDB(cfg)
	if aMysql.MysqlError != nil {
		endFunc(errors.New(fmt.Sprintf("连接数据库失败==%s", aMysql.MysqlError)))
		return
	}
	fmt.Printf("数据库==%s,连接成功", dbName)
	endFunc(nil)
}

// Close 关闭当前句柄的连接池。默认单例一般随进程生命周期存在,主要给 New 出来的实例用。
func (aMysql *GTORMMysql) Close() error {
	aMysql.mux.Lock()
	defer aMysql.mux.Unlock()

	if aMysql.MysqlDB == nil {
		return nil
	}
	sqlDB, err := aMysql.MysqlDB.DB()
	if err != nil {
		return err
	}
	if err := sqlDB.Close(); err != nil {
		return err
	}
	aMysql.MysqlDB = nil
	return nil
}

// mysqlDB 取当前句柄的 gorm.DB。只用读锁护住指针本身(与 OpenXxx / Close 的写互斥),
// 取到后即释放锁再跑 SQL——gorm.DB 自身并发安全,并发控制交给连接池。
func (aMysql *GTORMMysql) mysqlDB() (*gorm.DB, error) {
	aMysql.mux.RLock()
	defer aMysql.mux.RUnlock()

	if aMysql.MysqlDB == nil {
		return nil, fmt.Errorf("gtbox_orm_mysql: instance not opened")
	}
	return aMysql.MysqlDB, nil
}

// InsertData 数据不存在时插入。不持写锁跑 SQL:原先整段包在 mux.Lock 里,
// 全实例插入被串成一条,连接池开多大都只有 1 并发。
func (aMysql *GTORMMysql) InsertData(dataModel interface{}) error {
	db, err := aMysql.mysqlDB()
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
