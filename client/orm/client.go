package orm

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/pinguo/pgo2/core"
	"github.com/pinguo/pgo2/util"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"gorm.io/plugin/dbresolver"
)

// DB client component, wrapper for database/sql,
// support read-write splitting, configuration:
// components:
//      db:
//          driver: "mysql"
//          dsn:    "user:pass@tcp(127.0.0.1:3306)/db?charset=utf8&timeout=0.5s"
//          slaves: ["slave1 dsn", "slave2 dsn"]
//          maxIdleConn: 5
//          maxOpenConn: 10
//          sqlLog:false
//          maxConnTime: "1h"
//          slowLogTime: "100ms"
// mysql config: https://gorm.io/zh_CN/docs/connecting_to_the_database.html#MySQL
//          defaultStringSize:256
//          disableDatetimePrecision:false
//          dontSupportRenameIndex:false
//          dontSupportRenameColumn:false
//          skipInitializeWithVersion:false
// gorm config: https://gorm.io/zh_CN/docs/gorm_config.html
//          skipDefaultTransaction:false
//	        prepareStmt:false
//	        dryRun  :false
//	        allowGlobalUpdate :false
//	        singularTable :false
//	        tablePrefix :false
//	        disableAutomaticPing :false
//	        disableForeignKeyConstraintWhenMigrating :false
// log config:https://gorm.io/zh_CN/docs/logger.html
//          logLevel:1  //  Silent:1 Error:2  Warn:3 Info:4

func New(config map[string]interface{}) (interface{}, error) {
	c := &Client{}
	c.maxIdleConn = 5
	c.maxOpenConn = 0
	c.maxConnTime = time.Hour
	c.slowLogTime = 100 * time.Millisecond
	c.sqlLog = false

	if err := core.ClientConfigure(c, config); err != nil {
		return nil, err
	}

	if err := c.Init(); err != nil {
		return nil, err
	}

	return c, nil

}

type Client struct {
	Db *gorm.DB

	driver string   // driver name
	dsn    string   // master dsn
	slaves []string // slaves dsn

	// sql
	maxIdleConn int
	maxConnTime time.Duration
	maxOpenConn int

	// log
	slowLogTime time.Duration
	sqlLog      bool // Complete SQL log
	logLevel    logger.LogLevel

	// mysql
	defaultStringSize         uint
	disableDatetimePrecision  bool
	dontSupportRenameIndex    bool
	dontSupportRenameColumn   bool
	skipInitializeWithVersion bool

	// gorm
	skipDefaultTransaction                   bool
	prepareStmt                              bool
	dryRun                                   bool
	allowGlobalUpdate                        bool
	singularTable                            bool
	tablePrefix                              string
	disableAutomaticPing                     bool
	disableForeignKeyConstraintWhenMigrating bool

	dialect IFuncDialect

}

type IFuncDialect func(dsn string) gorm.Dialector

func (c *Client) Init() error {
	if c.driver == "" || c.dsn == "" {
		return errors.New("Db: driver and dsn are required")
	}

	if util.SliceSearchString(sql.Drivers(), c.driver) == -1 {
		return fmt.Errorf("Db: driver %s is not registered", c.driver)
	}

	db, err := gorm.Open(c.getDialect(), &gorm.Config{
		SkipDefaultTransaction:                   c.skipDefaultTransaction,
		PrepareStmt:                              c.prepareStmt,
		DryRun:                                   c.dryRun,
		AllowGlobalUpdate:                        c.allowGlobalUpdate,
		DisableAutomaticPing:                     c.disableAutomaticPing,
		DisableForeignKeyConstraintWhenMigrating: c.disableForeignKeyConstraintWhenMigrating,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   c.tablePrefix,
			SingularTable: c.singularTable,
		},
	})

	if err != nil {
		return fmt.Errorf("failed to connect database err: %s",err.Error())
	}

	if sqlDb, e := db.DB(); e != nil {
		return fmt.Errorf("Db: open %s error, %s", c.dsn, e.Error())
	} else {
		sqlDb.SetConnMaxLifetime(c.maxConnTime)
		sqlDb.SetMaxIdleConns(c.maxIdleConn)
		sqlDb.SetMaxOpenConns(c.maxOpenConn)
	}

	slaveLen := len(c.slaves)
	if slaveLen > 0 {
		replicas := make([]gorm.Dialector, slaveLen)
		for k, slaveStr := range c.slaves {
			replicas[k] = c.slaveDialect(slaveStr)
		}

		db.Use(dbresolver.Register(dbresolver.Config{
			Replicas: replicas,
		}))
	}

	c.Db = db

	return nil
}

func (c *Client) slaveDialect(slaveDsn string) gorm.Dialector{
	if c.dialect != nil {
		return c.dialect(c.dsn)
	}

	return mysql.Open(slaveDsn)
}

func (c *Client) getDialect() gorm.Dialector {
	if c.dialect != nil {
		return c.dialect(c.dsn)
	}

	return mysql.New(mysql.Config{
		DSN:                       c.dsn,
		DefaultStringSize:         c.defaultStringSize,
		DisableDatetimePrecision:  c.disableDatetimePrecision,
		DontSupportRenameIndex:    c.dontSupportRenameIndex,
		DontSupportRenameColumn:   c.dontSupportRenameColumn,
		SkipInitializeWithVersion: c.skipInitializeWithVersion,
	})
}

func (c *Client) SqlLog() bool {
	return c.sqlLog
}

func (c *Client) SlowLogTime() time.Duration {
	return c.slowLogTime
}

func (c *Client) LogLevel() logger.LogLevel {
	return c.logLevel
}

// SetDriver set driver db use, eg. "mysql"
func (c *Client) SetDriver(driver string) {
	c.driver = driver
}

// SetDsn set master dsn, the dsn is driver specified,
// eg. dsn format for github.com/go-sql-driver/mysql is
// [username[:password]@][protocol[(address)]]/dbname[?param=value]
func (c *Client) SetDsn(dsn string) {
	c.dsn = dsn
}

// SetMaxIdleConn set max idle conn, default is 5
func (c *Client) SetMaxIdleConn(maxIdleConn int) {
	c.maxIdleConn = maxIdleConn
}

// SetMaxOpenConn set max open conn, default is 0
func (c *Client) SetMaxOpenConn(maxOpenConn int) {
	c.maxOpenConn = maxOpenConn
}

// SetMaxConnTime set conn life time, default is 1h
func (c *Client) SetMaxConnTime(v string) error {
	if maxConnTime, err := time.ParseDuration(v); err != nil {
		return errors.New("Db.SetMaxConnTime error, " + err.Error())
	} else {
		c.maxConnTime = maxConnTime
	}

	return nil
}

// SetSlowTime set slow log time, default is 100ms
func (c *Client) SetSlowLogTime(v string) error {
	if slowLogTime, err := time.ParseDuration(v); err != nil {
		return errors.New("Db.SetSlowLogTime error, " + err.Error())
	} else {
		c.slowLogTime = slowLogTime
	}

	return nil
}

// SetSkipDefaultTransaction
// 对于写操作（创建、更新、删除），为了确保数据的完整性，GORM 会将它们封装在事务内运行。但这会降低性能，你可以在初始化时禁用这种方式
func (c *Client) SetSkipDefaultTransaction(v bool) {
	c.skipDefaultTransaction = v
}

// SetPrepareStmt
// 执行任何 SQL 时都创建并缓存预编译语句，可以提高后续的调用速度
func (c *Client) SetPrepareStmt(v bool) {
	c.prepareStmt = v
}

func (c *Client) SetSqlLog(sqlLog bool) {
	c.sqlLog = sqlLog
}

// string 类型字段的默认长度
func (c *Client) SetDefaultStringSize(v uint) {
	c.defaultStringSize = v
}

// 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
func (c *Client) SetDisableDatetimePrecision(v bool) {
	c.disableDatetimePrecision = v
}

// 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
func (c *Client) SetDontSupportRenameIndex(v bool) {
	c.dontSupportRenameIndex = v
}

// 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
func (c *Client) SetDontSupportRenameColumn(v bool) {
	c.dontSupportRenameColumn = v
}

// 根据当前 MySQL 版本自动配置
func (c *Client) SetSkipInitializeWithVersion(v bool) {
	c.skipInitializeWithVersion = v
}

// 生成 SQL 但不执行，可以用于准备或测试生成的 SQL，参考 会话 获取详情
func (c *Client) SetDryRun(v bool) {
	c.dryRun = v
}

// 启用全局 update/delete，查看 Session 获取详情
func (c *Client) SetAllowGlobalUpdate(v bool) {
	c.allowGlobalUpdate = v
}

// 使用单数表名，启用该选项
func (c *Client) SetSingularTable(v bool) {
	c.singularTable = v
}

// 表名前缀，`User` 的表名应该是 v+`_user`
func (c *Client) SetTablePrefix(v string) {
	c.tablePrefix = v
}

// 在完成初始化后，GORM 会自动 ping 数据库以检查数据库的可用性，若要禁用该特性，可将其设置为 true
func (c *Client) SetDisableAutomaticPing(v bool) {
	c.disableAutomaticPing = v
}

// 在 AutoMigrate 或 CreateTable 时，GORM 会自动创建外键约束，若要禁用该特性，可将其设置为 true，参考 迁移 获取详情。
func (c *Client) SetDisableForeignKeyConstraintWhenMigrating(v bool) {
	c.disableForeignKeyConstraintWhenMigrating = v
}

func (c *Client) SetLogLevel(v logger.LogLevel) {
	c.logLevel = v
}

// 设置数据库对象
func (c *Client) SetDialect(v IFuncDialect) {
	c.dialect = v
}

// SetSlaves set dsn for slaves
func (c *Client) SetSlaves(v []interface{}) {
	c.slaves = make([]string,0,len(v))
	for _, vv := range v {
		c.slaves = append(c.slaves, vv.(string))
	}
}