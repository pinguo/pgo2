package adapter

import (
	"context"
	"database/sql"
	"strings"

	"github.com/pinguo/pgo2"
	"github.com/pinguo/pgo2/client/orm"
	"github.com/pinguo/pgo2/iface"
	"github.com/pinguo/pgo2/logs"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

var OrmClass string
var OrmLogWriterClass string

func init() {
	container := pgo2.App().Container()
	OrmClass = container.Bind(&Orm{})
	OrmLogWriterClass = container.Bind(&LogWriter{})
}

// NewDb of db Client, add context support.
// usage: db := this.GetObj(adapter.NewOrm(ctr)).(*adapter.Orm/adapter.IOrm)
func NewOrm(ctr iface.IContext, componentId ...string) *Orm {
	id := DefaultOrmId
	if len(componentId) > 0 {
		id = componentId[0]
	}
	o := &Orm{}
	o.client = pgo2.App().Component(id, orm.New).(*orm.Client)
	o.componentId = id
	o.DB = o.dbSession(ctr)

	return o
}

type Orm struct {
	pgo2.Object
	client *orm.Client
	*gorm.DB
	componentId string
}

// NewDb of db Client, add context support.
// usage: db := this.GetObjBox(adapter.OrmClass).(*adapter.Orm/adapter.IOrm)
func (o *Orm) Prepare(componentId ...interface{}) {
	id := DefaultOrmId
	if len(componentId) > 0 {
		id = componentId[0].(string)
	}
	o.componentId = id
	o.client = pgo2.App().Component(id, orm.New).(*orm.Client)
	if len(componentId) < 2 {
		o.DB = o.dbSession(o.Context())
	}

}

func (o *Orm) defaultLogger(ctr iface.IContext) logger.Interface {
	return logger.New(
		o.GetObjBoxCtx(ctr, OrmLogWriterClass).(*LogWriter), // io writer
		logger.Config{
			SlowThreshold: o.client.SlowLogTime(), // 慢 SQL 阈值
			LogLevel:      o.client.LogLevel(),    // Log level
			Colorful:      false,                  // 禁用彩色打印
		},
	)
}

// new session
func (o *Orm) dbSession(ctr iface.IContext) *gorm.DB {
	return o.client.Db.Session(&gorm.Session{Logger: o.defaultLogger(ctr)})
}

func (o *Orm) clone(db *gorm.DB) IOrm {
	cO := o.GetObjBoxCtx(o.Context(), OrmClass, o.componentId, "notInitDbSession").(*Orm)
	cO.DB = db
	return cO
}

func (o *Orm) Assign(attrs ...interface{}) (tx IOrm) {
	return o.clone(o.DB.Assign(attrs...))
}

func (o *Orm) Attrs(attrs ...interface{}) (tx IOrm) {
	return o.clone(o.DB.Attrs(attrs...))
}

func (o *Orm) Begin(opts ...*sql.TxOptions) (tx IOrm) {
	return o.clone(o.DB.Begin(opts...))
}

func (o *Orm) Clauses(conds ...clause.Expression) (tx IOrm) {
	return o.clone(o.DB.Clauses(conds...))
}

func (o *Orm) Commit() (tx IOrm) {
	return o.clone(o.DB.Commit())
}

func (o *Orm) Count(count *int64) (tx IOrm) {
	profileKey := "orm.Count"
	o.Context().ProfileStart(profileKey)
	defer o.Context().ProfileStop(profileKey)
	return o.clone(o.DB.Count(count))
}

func (o *Orm) Create(value interface{}) (tx IOrm) {
	profileKey := "orm.Create"
	o.Context().ProfileStart(profileKey)
	defer o.Context().ProfileStop(profileKey)
	return o.clone(o.DB.Create(value))
}

func (o *Orm) CreateInBatches(value interface{}, batchSize int) (tx IOrm) {
	profileKey := "orm.CreateInBatches"
	o.Context().ProfileStart(profileKey)
	defer o.Context().ProfileStop(profileKey)
	return o.clone(o.DB.CreateInBatches(value, batchSize))
}

func (o *Orm) Debug() (tx IOrm) {
	return o.clone(o.DB.Debug())
}

func (o *Orm) Delete(value interface{}, conds ...interface{}) (tx IOrm) {
	profileKey := "orm.Delete"
	o.Context().ProfileStart(profileKey)
	defer o.Context().ProfileStop(profileKey)
	return o.clone(o.DB.Delete(value, conds...))
}

func (o *Orm) Distinct(args ...interface{}) (tx IOrm) {
	profileKey := "orm.Distinct"
	o.Context().ProfileStart(profileKey)
	defer o.Context().ProfileStop(profileKey)
	return o.clone(o.DB.Distinct(args...))
}

func (o *Orm) Exec(sql string, values ...interface{}) (tx IOrm) {
	profileKey := "orm.Exec"
	o.Context().ProfileStart(profileKey)
	defer o.Context().ProfileStop(profileKey)
	return o.clone(o.DB.Exec(sql, values...))
}

func (o *Orm) Find(dest interface{}, conds ...interface{}) (tx IOrm) {
	profileKey := "orm.Find"
	o.Context().ProfileStart(profileKey)
	defer o.Context().ProfileStop(profileKey)
	return o.clone(o.DB.Find(dest, conds...))
}

func (o *Orm) FindInBatches(dest interface{}, batchSize int, fc func(pTx *gorm.DB, batch int) error) (tx IOrm) {
	profileKey := "orm.FindInBatches"
	o.Context().ProfileStart(profileKey)
	defer o.Context().ProfileStop(profileKey)
	return o.clone(o.DB.FindInBatches(dest, batchSize, fc))
}

func (o *Orm) First(dest interface{}, conds ...interface{}) (tx IOrm) {
	profileKey := "orm.First"
	o.Context().ProfileStart(profileKey)
	defer o.Context().ProfileStop(profileKey)
	return o.clone(o.DB.First(dest, conds...))
}

func (o *Orm) FirstOrCreate(dest interface{}, conds ...interface{}) (tx IOrm) {
	profileKey := "orm.FirstOrCreate"
	o.Context().ProfileStart(profileKey)
	defer o.Context().ProfileStop(profileKey)
	return o.clone(o.DB.FirstOrCreate(dest, conds...))
}

func (o *Orm) FirstOrInit(dest interface{}, conds ...interface{}) (tx IOrm) {
	profileKey := "orm.FirstOrInit"
	o.Context().ProfileStart(profileKey)
	defer o.Context().ProfileStop(profileKey)
	return o.clone(o.DB.FirstOrInit(dest, conds...))
}

func (o *Orm) Group(name string) (tx IOrm) {
	profileKey := "orm.Group"
	o.Context().ProfileStart(profileKey)
	defer o.Context().ProfileStop(profileKey)
	return o.clone(o.DB.Group(name))
}

func (o *Orm) Having(query interface{}, args ...interface{}) (tx IOrm) {
	profileKey := "orm.Having"
	o.Context().ProfileStart(profileKey)
	defer o.Context().ProfileStop(profileKey)
	return o.clone(o.DB.Having(query, args...))
}

func (o *Orm) Joins(query string, args ...interface{}) (tx IOrm) {
	profileKey := "orm.Joins"
	o.Context().ProfileStart(profileKey)
	defer o.Context().ProfileStop(profileKey)
	return o.clone(o.DB.Joins(query, args...))
}

func (o *Orm) Last(dest string, conds ...interface{}) (tx IOrm) {
	profileKey := "orm.Last"
	o.Context().ProfileStart(profileKey)
	defer o.Context().ProfileStop(profileKey)
	return o.clone(o.DB.Last(dest, conds...))
}

func (o *Orm) Limit(limit int) (tx IOrm) {

	return o.clone(o.DB.Limit(limit))
}

func (o *Orm) Model(value interface{}) (tx IOrm) {
	return o.clone(o.DB.Model(value))
}

func (o *Orm) Not(query interface{}, args ...interface{}) (tx IOrm) {
	return o.clone(o.DB.Not(query, args...))
}

func (o *Orm) Offset(offset int) (tx IOrm) {
	return o.clone(o.DB.Offset(offset))
}

func (o *Orm) Omit(columns ...string) (tx IOrm) {
	return o.clone(o.DB.Omit(columns...))
}

func (o *Orm) Or(query interface{}, args ...interface{}) (tx IOrm) {
	return o.clone(o.DB.Or(query, args...))
}

func (o *Orm) Order(value interface{}) (tx IOrm) {
	return o.clone(o.DB.Order(value))
}

func (o *Orm) Pluck(column string, dest interface{}) (tx IOrm) {
	return o.clone(o.DB.Pluck(column, dest))
}

func (o *Orm) Preload(query string, args ...interface{}) (tx IOrm) {
	return o.clone(o.DB.Preload(query, args...))
}

func (o *Orm) Raw(sql string, values ...interface{}) (tx IOrm) {
	return o.clone(o.DB.Raw(sql, values...))
}

func (o *Orm) Rollback() (tx IOrm) {
	return o.clone(o.DB.Rollback())
}

func (o *Orm) RollbackTo(name string) (tx IOrm) {
	return o.clone(o.DB.RollbackTo(name))
}

func (o *Orm) Row() *sql.Row {
	profileKey := "orm.Row"
	o.Context().ProfileStart(profileKey)
	defer o.Context().ProfileStop(profileKey)
	return o.DB.Row()
}

func (o *Orm) Rows() (*sql.Rows, error) {
	profileKey := "orm.Rows"
	o.Context().ProfileStart(profileKey)
	defer o.Context().ProfileStop(profileKey)
	return o.DB.Rows()
}

func (o *Orm) Save(value interface{}) (tx IOrm) {
	profileKey := "orm.Save"
	o.Context().ProfileStart(profileKey)
	defer o.Context().ProfileStop(profileKey)
	return o.clone(o.DB.Save(value))
}

func (o *Orm) SavePoint(name string) (tx IOrm) {
	profileKey := "orm.SavePoint"
	o.Context().ProfileStart(profileKey)
	defer o.Context().ProfileStop(profileKey)
	return o.clone(o.DB.SavePoint(name))
}

func (o *Orm) Scan(dest interface{}) (tx IOrm) {
	profileKey := "orm.Scan"
	o.Context().ProfileStart(profileKey)
	defer o.Context().ProfileStop(profileKey)
	return o.clone(o.DB.Scan(dest))
}

func (o *Orm) Scopes(funcs ...func(*gorm.DB) *gorm.DB) (tx IOrm) {
	profileKey := "orm.Scopes"
	o.Context().ProfileStart(profileKey)
	defer o.Context().ProfileStop(profileKey)
	return o.clone(o.DB.Scopes(funcs...))
}

func (o *Orm) Select(query interface{}, args ...interface{}) (tx IOrm) {
	return o.clone(o.DB.Select(query, args...))
}

func (o *Orm) Session(config *gorm.Session) (tx IOrm) {
	if config == nil {
		config = &gorm.Session{}
	}

	if config.Logger == nil {
		config.Logger = o.defaultLogger(o.Context())
	}

	return o.clone(o.DB.Session(config))
}

func (o *Orm) InstanceSet(key string, value interface{}) (tx IOrm) {
	return o.clone(o.DB.InstanceSet(key, value))
}

func (o *Orm) Set(key string, value interface{}) (tx IOrm) {
	return o.clone(o.DB.Set(key, value))
}

func (o *Orm) Table(name string, args ...interface{}) (tx IOrm) {
	return o.clone(o.DB.Table(name, args...))
}

func (o *Orm) Take(dest interface{}, conds ...interface{}) (tx IOrm) {
	profileKey := "orm.Take"
	o.Context().ProfileStart(profileKey)
	defer o.Context().ProfileStop(profileKey)
	return o.clone(o.DB.Take(dest, conds...))
}

func (o *Orm) Unscoped() (tx IOrm) {
	return o.clone(o.DB.Unscoped())
}

func (o *Orm) Update(column string, value interface{}) (tx IOrm) {
	profileKey := "orm.Update"
	o.Context().ProfileStart(profileKey)
	defer o.Context().ProfileStop(profileKey)
	return o.clone(o.DB.Update(column, value))
}

func (o *Orm) UpdateColumn(column string, value interface{}) (tx IOrm) {
	profileKey := "orm.UpdateColumn"
	o.Context().ProfileStart(profileKey)
	defer o.Context().ProfileStop(profileKey)
	return o.clone(o.DB.UpdateColumn(column, value))
}

func (o *Orm) UpdateColumns(values interface{}) (tx IOrm) {
	profileKey := "orm.UpdateColumns"
	o.Context().ProfileStart(profileKey)
	defer o.Context().ProfileStop(profileKey)
	return o.clone(o.DB.UpdateColumns(values))
}

func (o *Orm) Updates(values interface{}) (tx IOrm) {
	profileKey := "orm.Updates"
	o.Context().ProfileStart(profileKey)
	defer o.Context().ProfileStop(profileKey)
	return o.clone(o.DB.Updates(values))
}

func (o *Orm) Where(query interface{}, args ...interface{}) (tx IOrm) {

	return o.clone(o.DB.Where(query, args...))
}

func (o *Orm) WithContext(ctx context.Context) (tx IOrm) {
	o.Callback()
	return o.clone(o.DB.WithContext(ctx))
}

func (o *Orm) GetError() error {
	return o.Error
}

func (o *Orm) GetRowsAffected() int64 {
	return o.RowsAffected
}

func (o *Orm) GetStatement() *gorm.Statement {
	return o.Statement
}

func (o *Orm) GetConfig() *gorm.Config {
	return o.Config
}

func (o *Orm) SqlDB() (*sql.DB, error) {
	return o.DB.DB()
}

type LogWriter struct {
	pgo2.Object
	logger logs.ILogger
}

func (l *LogWriter) Prepare() {
	l.logger = l.Context()
}

func (l *LogWriter) Printf(msg string, data ...interface{}) {
	if strings.Index(msg, "[warn]") >= 0 {
		l.logger.Warn(msg, data...)
	} else if strings.Index(msg, "[error]") >= 0 {
		l.logger.Error(msg, data...)
	} else {
		l.logger.Info(msg, data...)
	}

}
