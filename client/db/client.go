package db

import (
    "database/sql"
    "errors"
    "fmt"
    "time"

    "github.com/pinguo/pgo2/core"
    "github.com/pinguo/pgo2/util"
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
func New(config map[string]interface{}) (interface{}, error) {
    c := &Client{}
    c.maxIdleConn = 5
    c.maxOpenConn = 0
    c.maxConnTime = time.Hour
    c.slowLogTime = 100 * time.Millisecond
    c.sqlLog = false

    if err := core.ClientConfigure(c, config);err!=nil{
        return nil,err
    }

    if err := c.Init();err!=nil{
        return nil,err
    }

    return c,nil

}

type Client struct {
    driver string   // driver name
    dsn    string   // master dsn
    slaves []string // slaves dsn

    maxIdleConn int
    maxConnTime time.Duration
    slowLogTime time.Duration
    maxOpenConn int

    sqlLog bool // Complete SQL log

    masterDb *sql.DB   // master db instance
    slaveDbs []*sql.DB // slave db instances
}

func (c *Client) Init() error{
    if c.driver == "" || c.dsn == "" {
        return errors.New("Db: driver and dsn are required")
    }

    if util.SliceSearchString(sql.Drivers(), c.driver) == -1 {
        return fmt.Errorf("Db: driver %s is not registered", c.driver)
    }

    // create master db instance
    if db, e := sql.Open(c.driver, c.dsn); e != nil {
        return fmt.Errorf("Db: open %s error, %s", c.dsn, e.Error())
    } else {
        db.SetConnMaxLifetime(c.maxConnTime)
        db.SetMaxIdleConns(c.maxIdleConn)
        db.SetMaxOpenConns(c.maxOpenConn)
        c.masterDb = db
    }

    // create slave db instances
    for _, dsn := range c.slaves {
        if db, e := sql.Open(c.driver, dsn); e != nil {
            return fmt.Errorf("Db: open %s error, %s", dsn, e.Error())
        } else {
            db.SetConnMaxLifetime(c.maxConnTime)
            db.SetMaxIdleConns(c.maxIdleConn)
            db.SetMaxOpenConns(c.maxOpenConn)
            c.slaveDbs = append(c.slaveDbs, db)
        }
    }

    return nil
}

func (c *Client) SqlLog() bool{
    return c.sqlLog
}

func (c *Client) SlowLogTime() time.Duration{
    return c.slowLogTime
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

// SetSlaves set dsn for slaves
func (c *Client) SetSlaves(v []interface{}) {
    for _, vv := range v {
        c.slaves = append(c.slaves, vv.(string))
    }
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
func (c *Client) SetMaxConnTime(v string) error{
    if maxConnTime, err := time.ParseDuration(v); err != nil {
        return errors.New("Db.SetMaxConnTime error, " + err.Error())
    } else {
        c.maxConnTime = maxConnTime
    }

    return nil
}

// SetSlowTime set slow log time, default is 100ms
func (c *Client) SetSlowLogTime(v string) error{
    if slowLogTime, err := time.ParseDuration(v); err != nil {
        return errors.New("Db.SetSlowLogTime error, " + err.Error())
    } else {
        c.slowLogTime = slowLogTime
    }

    return nil
}

func (c *Client) SetSqlLog(sqlLog bool) {
    c.sqlLog = sqlLog
}

// GetDb get a master or slave db instance
func (c *Client) GetDb(master bool) *sql.DB {
    if num := len(c.slaveDbs); !master && num > 0 {
        idx := 0
        if num > 1 {
            idx = (time.Now().Nanosecond() / 1000) % num
        }

        return c.slaveDbs[idx]
    }

    return c.masterDb
}
