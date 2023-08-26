package gtbox_orm_sqlite

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
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

func Instance() *GTORMSqlite {
	gtSqliteOnce.Do(func() {
		sqliteInstance = &GTORMSqlite{}
	})
	return sqliteInstance
}

func (gtSqlite *GTORMSqlite) OpenSqlite(sqlitePath string) {
	gtSqlite.mux.Lock()
	defer gtSqlite.mux.Unlock()
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	gtSqlite.SqliteDB, gtSqlite.SqliteError = gorm.Open(sqlite.Open(sqlitePath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if gtSqlite.SqliteError != nil {
		println("连接数据库失败==", gtSqlite.SqliteDB, gtSqlite.SqliteError)
	}
}

func (gtSqlite *GTORMSqlite) InsertData(dataModel interface{}) error {
	gtSqlite.mux.Lock()
	defer gtSqlite.mux.Unlock()

	result := gtSqlite.SqliteDB.Where(dataModel).Limit(1).Find(dataModel)
	err := result.Error
	if err != nil {
		return err
	}
	if result.RowsAffected == 0 {
		if cuerr := gtSqlite.SqliteDB.Create(dataModel).Error; cuerr != nil {
			return cuerr
		}
	}
	return nil
}

func (gtSqlite *GTORMSqlite) QueryData(dataModel interface{}, conditions ...interface{}) error {
	gtSqlite.mux.RLock()
	defer gtSqlite.mux.RUnlock()

	if err := gtSqlite.SqliteDB.Where(dataModel, conditions...).Find(dataModel).Error; err != nil {
		return err
	}
	return nil
}
