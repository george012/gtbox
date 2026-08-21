package gtbox_orm_mysql

import (
	"fmt"
	"time"

	"github.com/george012/gtbox/gtbox_orm/gtbox_orm_config"
	sqldriver "github.com/go-sql-driver/mysql"
)

// 连接池默认值，与 OPenMysql 原有硬编码一致；GTORMMysqlConfig 对应字段留零值即取这组值。
const (
	defaultMysqlConnMaxOpen     = 5
	defaultMysqlConnMaxIdle     = 2
	defaultMysqlConnMaxIdleTime = time.Minute
)

// GTORMMysqlConfig 一个 mysql 数据源的连接身份与连接池参数。
// 连接身份字段(DBUser/DBName/DBHost/DBPort)无默认值，零值直接报错——公用层不兜底 localhost/3306。
// 池参数零值 = 取上面那组默认值；ConnMaxLifetime 零值 = 不限制连接存活时长(database/sql 默认)。
type GTORMMysqlConfig struct {
	DBUser          string
	DBPwd           string // 允许为空，对应无密码账号
	DBName          string
	DBHost          string
	DBPort          int
	DBTimeZone      gtbox_orm_config.GTORMTimeZone
	ConnMaxOpen     int
	ConnMaxIdle     int
	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration
}

// validate 参数校验；ConnMaxIdle > ConnMaxOpen 时 database/sql 会静默把 idle 压到 open，
// 这里改为直接拒绝，避免配置写错却查不出来。
func (cfg *GTORMMysqlConfig) validate() error {
	if cfg.DBUser == "" {
		return fmt.Errorf("gtbox_orm_mysql: DBUser is required")
	}
	if cfg.DBName == "" {
		return fmt.Errorf("gtbox_orm_mysql: DBName is required")
	}
	if cfg.DBHost == "" {
		return fmt.Errorf("gtbox_orm_mysql: DBHost is required")
	}
	if cfg.DBPort <= 0 || cfg.DBPort > 65535 {
		return fmt.Errorf("gtbox_orm_mysql: DBPort must be in 1-65535, got %d", cfg.DBPort)
	}
	if cfg.ConnMaxOpen < 0 {
		return fmt.Errorf("gtbox_orm_mysql: ConnMaxOpen must be >= 0, got %d", cfg.ConnMaxOpen)
	}
	if cfg.ConnMaxIdle < 0 {
		return fmt.Errorf("gtbox_orm_mysql: ConnMaxIdle must be >= 0, got %d", cfg.ConnMaxIdle)
	}
	if cfg.ConnMaxIdleTime < 0 {
		return fmt.Errorf("gtbox_orm_mysql: ConnMaxIdleTime must be >= 0, got %s", cfg.ConnMaxIdleTime)
	}
	if cfg.ConnMaxLifetime < 0 {
		return fmt.Errorf("gtbox_orm_mysql: ConnMaxLifetime must be >= 0, got %s", cfg.ConnMaxLifetime)
	}
	if cfg.ConnMaxOpen > 0 && cfg.ConnMaxIdle > cfg.ConnMaxOpen {
		return fmt.Errorf("gtbox_orm_mysql: ConnMaxIdle(%d) must be <= ConnMaxOpen(%d)", cfg.ConnMaxIdle, cfg.ConnMaxOpen)
	}
	return nil
}

// poolParams 返回补齐默认值后的池参数
func (cfg *GTORMMysqlConfig) poolParams() (maxOpen int, maxIdle int, idleTime time.Duration) {
	maxOpen, maxIdle, idleTime = cfg.ConnMaxOpen, cfg.ConnMaxIdle, cfg.ConnMaxIdleTime
	if maxOpen == 0 {
		maxOpen = defaultMysqlConnMaxOpen
	}
	if maxIdle == 0 {
		maxIdle = defaultMysqlConnMaxIdle
	}
	if idleTime == 0 {
		idleTime = defaultMysqlConnMaxIdleTime
	}
	return maxOpen, maxIdle, idleTime
}

// mysqlTimeZoneLocation GTORMTimeZone 到 *time.Location 的映射。
// GTORMTimeZone.String() 产出的是 DSN 转义形式(Asia%2FShanghai)，不能直接喂 time.LoadLocation。
// 宿主缺 tzdata 时返回 error，不静默退回 UTC——时区错位比连不上更难查。
func mysqlTimeZoneLocation(timeZone gtbox_orm_config.GTORMTimeZone) (*time.Location, error) {
	switch timeZone {
	case gtbox_orm_config.GTORMTimeZoneShangHai:
		return time.LoadLocation("Asia/Shanghai")
	default:
		return time.UTC, nil
	}
}

// buildMysqlDSN 由 driver 自己的 Config.FormatDSN 生成 DSN。
// 密码含 @ : / 等字符不需要转义——driver 的 ParseDSN 按「最后一个 @」「最后一个 /」切分,
// 用 Config 是为了参数编码归 driver 管:loc 由 Loc 字段带 URL 转义写出,
// 不依赖 GTORMTimeZone.String() 里那份手写的 %2F。
func buildMysqlDSN(cfg *GTORMMysqlConfig, loc *time.Location) string {
	dsnCfg := sqldriver.NewConfig()
	dsnCfg.User = cfg.DBUser
	dsnCfg.Passwd = cfg.DBPwd
	dsnCfg.Net = "tcp"
	dsnCfg.Addr = fmt.Sprintf("%s:%d", cfg.DBHost, cfg.DBPort)
	dsnCfg.DBName = cfg.DBName
	dsnCfg.ParseTime = true
	dsnCfg.Loc = loc
	return dsnCfg.FormatDSN()
}
