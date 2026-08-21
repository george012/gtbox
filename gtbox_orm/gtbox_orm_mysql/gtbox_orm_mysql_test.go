package gtbox_orm_mysql

import (
	"strings"
	"testing"
	"time"

	"github.com/george012/gtbox/gtbox_orm/gtbox_orm_config"
	sqldriver "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// TestGTORMMysqlConfigValidate 连接身份缺失、池参数非法必须显式拒绝,不兜底默认值
func TestGTORMMysqlConfigValidate(t *testing.T) {
	base := GTORMMysqlConfig{User: "root", Pwd: "pwd", DBName: "shop", Host: "127.0.0.1", Port: 3306}

	cases := []struct {
		name    string
		mutate  func(cfg *GTORMMysqlConfig)
		wantErr bool
	}{
		{"合法配置", func(cfg *GTORMMysqlConfig) {}, false},
		{"空密码合法", func(cfg *GTORMMysqlConfig) { cfg.Pwd = "" }, false},
		{"缺 User", func(cfg *GTORMMysqlConfig) { cfg.User = "" }, true},
		{"缺 DBName", func(cfg *GTORMMysqlConfig) { cfg.DBName = "" }, true},
		{"缺 Host", func(cfg *GTORMMysqlConfig) { cfg.Host = "" }, true},
		{"Port 为 0", func(cfg *GTORMMysqlConfig) { cfg.Port = 0 }, true},
		{"Port 越界", func(cfg *GTORMMysqlConfig) { cfg.Port = 70000 }, true},
		{"MaxOpenConns 负数", func(cfg *GTORMMysqlConfig) { cfg.MaxOpenConns = -1 }, true},
		{"MaxIdleConns 超过 MaxOpenConns", func(cfg *GTORMMysqlConfig) { cfg.MaxOpenConns, cfg.MaxIdleConns = 4, 8 }, true},
		{"ConnMaxLifetime 负数", func(cfg *GTORMMysqlConfig) { cfg.ConnMaxLifetime = -time.Second }, true},
	}

	for _, aCase := range cases {
		t.Run(aCase.name, func(t *testing.T) {
			cfg := base
			aCase.mutate(&cfg)
			err := cfg.validate()
			if (err != nil) != aCase.wantErr {
				t.Fatalf("validate() err=%v, wantErr=%v", err, aCase.wantErr)
			}
			// New 在校验不过时必须先返回 error,不能进到真实拨号
			if aCase.wantErr {
				if _, newErr := New(cfg); newErr == nil {
					t.Fatalf("New() 应在参数校验阶段失败")
				}
			}
		})
	}
}

// TestGTORMMysqlConfigPoolParams 池参数零值取包内默认值,与 OPenMysql 原有硬编码一致
func TestGTORMMysqlConfigPoolParams(t *testing.T) {
	zeroCfg := GTORMMysqlConfig{}
	maxOpen, maxIdle, idleTime := zeroCfg.poolParams()
	if maxOpen != defaultMysqlMaxOpenConns || maxIdle != defaultMysqlMaxIdleConns || idleTime != defaultMysqlConnMaxIdleTime {
		t.Fatalf("零值池参数=%d/%d/%s, want %d/%d/%s", maxOpen, maxIdle, idleTime,
			defaultMysqlMaxOpenConns, defaultMysqlMaxIdleConns, defaultMysqlConnMaxIdleTime)
	}

	customCfg := GTORMMysqlConfig{MaxOpenConns: 32, MaxIdleConns: 8, ConnMaxIdleTime: 30 * time.Second}
	maxOpen, maxIdle, idleTime = customCfg.poolParams()
	if maxOpen != 32 || maxIdle != 8 || idleTime != 30*time.Second {
		t.Fatalf("自定义池参数被改写=%d/%d/%s", maxOpen, maxIdle, idleTime)
	}
}

// TestBuildMysqlDSN 生成的 DSN 必须能被 driver 自己解回原值(密码含 @ : / 也不例外),
// 且时区非 UTC 时 loc 参数带 URL 转义
func TestBuildMysqlDSN(t *testing.T) {
	cfg := GTORMMysqlConfig{User: "root", Pwd: "p@ss:w/ord", DBName: "shop", Host: "10.0.0.7", Port: 3307}
	loc, err := mysqlTimeZoneLocation(gtbox_orm_config.GTORMTimeZoneShangHai)
	if err != nil {
		t.Fatalf("mysqlTimeZoneLocation() err=%v", err)
	}
	dsn := buildMysqlDSN(&cfg, loc)

	if !strings.Contains(dsn, "parseTime=true") || !strings.Contains(dsn, "loc=Asia%2FShanghai") {
		t.Fatalf("DSN 缺 parseTime / loc: %s", dsn)
	}

	parsed, err := sqldriver.ParseDSN(dsn)
	if err != nil {
		t.Fatalf("driver 解不动自己生成的 DSN: %s, err=%v", dsn, err)
	}
	if parsed.User != cfg.User || parsed.Passwd != cfg.Pwd || parsed.DBName != cfg.DBName ||
		parsed.Addr != "10.0.0.7:3307" || parsed.Net != "tcp" || !parsed.ParseTime ||
		parsed.Loc.String() != "Asia/Shanghai" {
		t.Fatalf("DSN 解析结果与配置不符: %+v", parsed)
	}
}

// TestMysqlTimeZoneLocation 上海时区映射到真实 Location,不是 DSN 转义串
func TestMysqlTimeZoneLocation(t *testing.T) {
	loc, err := mysqlTimeZoneLocation(gtbox_orm_config.GTORMTimeZoneShangHai)
	if err != nil {
		t.Fatalf("mysqlTimeZoneLocation() err=%v", err)
	}
	if loc.String() != "Asia/Shanghai" {
		t.Fatalf("loc=%s, want Asia/Shanghai", loc.String())
	}
}

// TestGTORMMysqlOpenWithConfigRejectReopen 句柄已打开时重复 Open 必须显式拒绝,避免旧连接池悬空
func TestGTORMMysqlOpenWithConfigRejectReopen(t *testing.T) {
	aMysql := &GTORMMysql{MysqlDB: &gorm.DB{}}
	err := aMysql.OpenWithConfig(GTORMMysqlConfig{User: "root", DBName: "shop", Host: "127.0.0.1", Port: 3306})
	if err == nil || !strings.Contains(err.Error(), "already opened") {
		t.Fatalf("重复 Open err=%v, want already opened", err)
	}
}
