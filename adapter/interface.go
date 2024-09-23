package adapter

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/qiniu/qmgo"
	opts "github.com/qiniu/qmgo/options"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/pinguo/pgo2/client/es"
	"github.com/pinguo/pgo2/client/maxmind"
	"github.com/pinguo/pgo2/client/memcache"
	"github.com/pinguo/pgo2/client/mongodb"
	"github.com/pinguo/pgo2/client/phttp"
	"github.com/pinguo/pgo2/client/rabbitmq"
	"github.com/pinguo/pgo2/client/redis"
	"github.com/pinguo/pgo2/iface"
	"github.com/pinguo/pgo2/value"
)

type IMemory interface {
	SetPanicRecover(v bool)
	Get(key string) *value.Value
	MGet(keys []string) map[string]*value.Value
	Set(key string, value interface{}, expire ...time.Duration) bool
	MSet(items map[string]interface{}, expire ...time.Duration) bool
	Add(key string, value interface{}, expire ...time.Duration) bool
	MAdd(items map[string]interface{}, expire ...time.Duration) bool
	Del(key string) bool
	MDel(keys []string) bool
	Exists(key string) bool
	Incr(key string, delta int) int
}

type IHttp interface {
	SetPanicRecover(v bool)
	Get(addr string, data interface{}, option ...*phttp.Option) *http.Response
	Post(addr string, data interface{}, option ...*phttp.Option) *http.Response
	Do(req *http.Request, option ...*phttp.Option) *http.Response
	DoMulti(requests []*http.Request, option ...*phttp.Option) []*http.Response
}

type IMaxMind interface {
	GeoByIp(ip string, args ...interface{}) *maxmind.Geo
}

type IMemCache interface {
	SetPanicRecover(v bool)
	Get(key string) *value.Value
	MGet(keys []string) map[string]*value.Value
	Set(key string, value interface{}, expire ...time.Duration) bool
	MSet(items map[string]interface{}, expire ...time.Duration) bool
	Add(key string, value interface{}, expire ...time.Duration) bool
	MAdd(items map[string]interface{}, expire ...time.Duration) bool
	Del(key string) bool
	MDel(keys []string) bool
	Exists(key string) bool
	Incr(key string, delta int) int
	Retrieve(cmd, key string) *memcache.Item
	MultiRetrieve(cmd string, keys []string) []*memcache.Item
	Store(cmd string, item *memcache.Item, expire ...time.Duration) bool
	MultiStore(cmd string, items []*memcache.Item, expire ...time.Duration) bool
}

type IMongo interface {
	FindOne(query interface{}, result interface{}, options ...bson.M) error
	FindAll(query interface{}, result interface{}, options ...bson.M) error
	FindAndModify(query interface{}, change mgo.Change, result interface{}, options ...bson.M) error
	FindDistinct(query interface{}, key string, result interface{}, options ...bson.M) error
	InsertOne(doc interface{}) error
	InsertAll(docs []interface{}) error
	UpdateOne(query interface{}, update interface{}) error
	UpdateAll(query interface{}, update interface{}) error
	UpdateOrInsert(query interface{}, update interface{}) error
	DeleteOne(query interface{}) error
	DeleteAll(query interface{}) error
	Count(query interface{}, options ...bson.M) (int, error)
	PipeOne(pipeline interface{}, result interface{}) error
	PipeAll(pipeline interface{}, result interface{}) error
	MapReduce(query interface{}, job *mgo.MapReduce, result interface{}, options ...bson.M) error
}

type IRabbitMq interface {
	SetPanicRecover(v bool)
	ExchangeDeclare(dftExchange ...*rabbitmq.ExchangeData)
	Publish(opCode string, data interface{}, dftOpUid ...string) bool
	PublishExchange(serviceName, exchangeName, exchangeType, opCode string, data interface{}, dftOpUid ...string) bool
	ChannelBox() *rabbitmq.ChannelBox
	GetConsumeChannelBox(queueName string, opCodes []string, dftExchange ...*rabbitmq.ExchangeData) *rabbitmq.ChannelBox
	Consume(queueName string, opCodes []string, limit int, autoAck, noWait, exclusive bool) <-chan amqp.Delivery
	ConsumeExchange(exchangeName, exchangeType, queueName string, opCodes []string, limit int, autoAck, noWait, exclusive bool) <-chan amqp.Delivery
	DecodeBody(d amqp.Delivery, ret interface{}) error
	DecodeHeaders(d amqp.Delivery) *rabbitmq.RabbitHeaders
}

type IDb interface {
	SetMaster(v bool)
	GetDb(master bool) *sql.DB
	Begin(opts ...*sql.TxOptions) ITx
	BeginContext(ctx context.Context, opts *sql.TxOptions) ITx
	QueryOne(query string, args ...interface{}) IRow
	QueryOneContext(ctx context.Context, query string, args ...interface{}) IRow
	Query(query string, args ...interface{}) *sql.Rows
	QueryContext(ctx context.Context, query string, args ...interface{}) *sql.Rows
	Exec(query string, args ...interface{}) sql.Result
	ExecContext(ctx context.Context, query string, args ...interface{}) sql.Result
	PrepareSql(query string) IStmt
	PrepareContext(ctx context.Context, query string) IStmt
}

type ITx interface {
	iface.IObject
	Commit() bool
	Rollback() bool
	QueryOne(query string, args ...interface{}) IRow
	QueryOneContext(ctx context.Context, query string, args ...interface{}) IRow
	Query(query string, args ...interface{}) *sql.Rows
	QueryContext(ctx context.Context, query string, args ...interface{}) *sql.Rows
	Exec(query string, args ...interface{}) sql.Result
	ExecContext(ctx context.Context, query string, args ...interface{}) sql.Result
	PrepareSql(query string) IStmt
	PrepareContext(ctx context.Context, query string) IStmt
}

type IRow interface {
	iface.IObject
	Scan(dest ...interface{}) error
}

type IStmt interface {
	iface.IObject
	Close()
	QueryOne(args ...interface{}) IRow
	QueryOneContext(ctx context.Context, args ...interface{}) IRow
	Query(args ...interface{}) *sql.Rows
	QueryContext(ctx context.Context, args ...interface{}) *sql.Rows
	Exec(args ...interface{}) sql.Result
	ExecContext(ctx context.Context, args ...interface{}) sql.Result
}

type IOrm interface {
	iface.IObject
	Assign(attrs ...interface{}) IOrm
	Attrs(attrs ...interface{}) IOrm
	Begin(opts ...*sql.TxOptions) IOrm
	Clauses(conds ...clause.Expression) IOrm
	Commit() IOrm
	Count(count *int64) IOrm
	Create(value interface{}) IOrm
	CreateInBatches(value interface{}, batchSize int) IOrm
	Debug() IOrm
	Delete(value interface{}, conds ...interface{}) IOrm
	Distinct(args ...interface{}) IOrm
	Exec(sql string, values ...interface{}) IOrm
	Find(dest interface{}, conds ...interface{}) IOrm
	FindInBatches(dest interface{}, batchSize int, fc func(pTx *gorm.DB, batch int) error) IOrm
	First(dest interface{}, conds ...interface{}) IOrm
	FirstOrCreate(dest interface{}, conds ...interface{}) IOrm
	FirstOrInit(dest interface{}, conds ...interface{}) IOrm
	Group(name string) IOrm
	Having(query interface{}, args ...interface{}) IOrm
	Joins(query string, args ...interface{}) IOrm
	Last(dest string, conds ...interface{}) IOrm
	Limit(limit int) IOrm
	Model(value interface{}) IOrm
	Not(query interface{}, args ...interface{}) IOrm
	Offset(offset int) IOrm
	Omit(columns ...string) IOrm
	Or(query interface{}, args ...interface{}) IOrm
	Order(value interface{}) IOrm
	Pluck(column string, dest interface{}) IOrm
	Preload(query string, args ...interface{}) IOrm
	Raw(sql string, values ...interface{}) IOrm
	Rollback() IOrm
	RollbackTo(name string) IOrm
	Row() *sql.Row
	Rows() (*sql.Rows, error)
	Save(value interface{}) IOrm
	SavePoint(name string) IOrm
	Scan(dest interface{}) IOrm
	Scopes(funcs ...func(*gorm.DB) *gorm.DB) IOrm
	Select(query interface{}, args ...interface{}) IOrm
	Session(config *gorm.Session) IOrm
	InstanceSet(key string, value interface{}) IOrm
	Set(key string, value interface{}) IOrm
	Table(name string, args ...interface{}) IOrm
	Take(dest interface{}, conds ...interface{}) IOrm
	Unscoped() IOrm
	Update(column string, value interface{}) IOrm
	UpdateColumn(column string, value interface{}) IOrm
	UpdateColumns(values interface{}) IOrm
	Updates(values interface{}) IOrm
	Where(query interface{}, args ...interface{}) IOrm
	WithContext(ctx context.Context) IOrm
	AddError(err error) error
	Association(column string) *gorm.Association
	AutoMigrate(dst ...interface{}) error
	SqlDB() (*sql.DB, error)
	Get(key string) (interface{}, bool)
	InstanceGet(key string) (interface{}, bool)
	ScanRows(rows *sql.Rows, dest interface{}) error
	SetupJoinTable(model interface{}, field string, joinTable interface{}) error
	Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error)
	Use(plugin gorm.Plugin) (err error)
	GetError() error
	GetRowsAffected() int64
	GetStatement() *gorm.Statement
	GetConfig() *gorm.Config
}

type IMongodb interface {
	iface.IObject
	Count(query interface{}) (int, error)
	FindAndModify(query interface{}, change qmgo.Change, result interface{}, options ...bson.M) error
	DeleteAll(query interface{}) error
	DeleteOne(query interface{}) error
	InsertAll(docs interface{}) error
	PipeAll(query interface{}, desc interface{}) error
	UpdateOrInsert(query interface{}, doc interface{}) error
	FindOne(query interface{}, doc interface{}, options ...bson.M) error
	FindAll(query interface{}, doc interface{}, options ...bson.M) error
	InsertOne(doc interface{}, opts ...opts.InsertOneOptions) (result *qmgo.InsertOneResult, err error)
	InsertOneCtx(ctx context.Context, doc interface{}, opts ...opts.InsertOneOptions) (result *qmgo.InsertOneResult, err error)
	InsertMany(docs interface{}, opts ...opts.InsertManyOptions) (result *qmgo.InsertManyResult, err error)
	InsertManyCtx(ctx context.Context, docs interface{}, opts ...opts.InsertManyOptions) (result *qmgo.InsertManyResult, err error)
	Upsert(filter interface{}, replacement interface{}, opts ...opts.UpsertOptions) (result *qmgo.UpdateResult, err error)
	UpsertCtx(ctx context.Context, filter interface{}, replacement interface{}, opts ...opts.UpsertOptions) (result *qmgo.UpdateResult, err error)
	UpsertId(id interface{}, replacement interface{}, opts ...opts.UpsertOptions) (result *qmgo.UpdateResult, err error)
	UpsertIdCtx(ctx context.Context, id interface{}, replacement interface{}, opts ...opts.UpsertOptions) (result *qmgo.UpdateResult, err error)
	UpdateOne(filter interface{}, update interface{}, opts ...opts.UpdateOptions) (err error)
	UpdateOneCtx(ctx context.Context, filter interface{}, update interface{}, opts ...opts.UpdateOptions) (err error)
	UpdateId(id interface{}, update interface{}, opts ...opts.UpdateOptions) (err error)
	UpdateIdCtx(ctx context.Context, id interface{}, update interface{}, opts ...opts.UpdateOptions) (err error)
	UpdateAll(filter interface{}, update interface{}, opts ...opts.UpdateOptions) (result *qmgo.UpdateResult, err error)
	UpdateAllCtx(ctx context.Context, filter interface{}, update interface{}, opts ...opts.UpdateOptions) (result *qmgo.UpdateResult, err error)
	ReplaceOne(filter interface{}, doc interface{}, opts ...opts.ReplaceOptions) (err error)
	ReplaceOneCtx(ctx context.Context, filter interface{}, doc interface{}, opts ...opts.ReplaceOptions) (err error)
	Remove(filter interface{}, opts ...opts.RemoveOptions) (err error)
	RemoveCtx(ctx context.Context, filter interface{}, opts ...opts.RemoveOptions) (err error)
	RemoveId(id interface{}, opts ...opts.RemoveOptions) (err error)
	RemoveIdCtx(ctx context.Context, id interface{}, opts ...opts.RemoveOptions) (err error)
	RemoveAll(filter interface{}, opts ...opts.RemoveOptions) (result *qmgo.DeleteResult, err error)
	RemoveAllCtx(ctx context.Context, filter interface{}, opts ...opts.RemoveOptions) (result *qmgo.DeleteResult, err error)
	Aggregate(pipeline interface{}) IMongodbAggregate
	AggregateCtx(ctx context.Context, pipeline interface{}) IMongodbAggregate
	EnsureIndexes(uniques []string, indexes []string) (err error)
	EnsureIndexesCtx(ctx context.Context, uniques []string, indexes []string) (err error)
	CreateIndexes(indexes []opts.IndexModel) (err error)
	CreateIndexesCtx(ctx context.Context, indexes []opts.IndexModel) (err error)
	CreateOneIndex(index opts.IndexModel) error
	CreateOneIndexCtx(ctx context.Context, index opts.IndexModel) (err error)
	DropAllIndexes() (err error)
	DropAllIndexesCtx(ctx context.Context) (err error)
	DropIndex(indexes []string) error
	DropIndexCtx(ctx context.Context, indexes []string) error
	DropCollection() error
	DropCollectionCtx(ctx context.Context) error
	CloneCollection() (*mongo.Collection, error)
	GetCollectionName() string
	Find(filter interface{}, options ...opts.FindOptions) IMongodbQuery
	FindCtx(ctx context.Context, filter interface{}, options ...opts.FindOptions) IMongodbQuery
	GetClient() *mongodb.Client
	Session() (*qmgo.Session, error)
	DoTransaction(ctx context.Context, callback func(sessCtx context.Context) (interface{}, error)) (interface{}, error)
}

type IMongodbQuery interface {
	qmgo.QueryI
}

type IMongodbAggregate interface {
	qmgo.AggregateI
}

type IRedis interface {
	SetPanicRecover(v bool)
	Get(key string) *value.Value
	MGet(keys []string) map[string]*value.Value
	Set(key string, value interface{}, expire ...time.Duration) bool
	MSet(items map[string]interface{}, expire ...time.Duration) bool
	Add(key string, value interface{}, expire ...time.Duration) bool
	MAdd(items map[string]interface{}, expire ...time.Duration) bool
	Del(key string) bool
	MDel(keys []string) bool
	Exists(key string) bool
	Incr(key string, delta int) int
	IncrBy(key string, delta int64) (int64, error)
	Do(cmd string, args ...interface{}) interface{}
	ExpireAt(key string, timestamp int64) bool
	Expire(key string, expire time.Duration) bool
	RPush(key string, values ...interface{}) bool
	LPush(key string, values ...interface{}) bool
	RPop(key string) *value.Value
	LPop(key string) *value.Value
	LLen(key string) int64
	HDel(key string, fields ...interface{}) int64
	HExists(key, field string) bool
	HSet(key string, fv ...interface{}) bool
	HGet(key, field string) *value.Value
	HGetAll(key string) map[string]*value.Value
	HMSet(key string, fv ...interface{}) bool
	HMGet(key string, fields ...interface{}) map[string]*value.Value
	HIncrBy(key, field string, delta int64) (int64, error)
	ZRange(key string, start, end int) []*value.Value
	ZRevRange(key string, start, end int) []*value.Value
	ZRangeWithScores(key string, start, end int) []*redis.ZV
	ZRevRangeWithScores(key string, start, end int) []*redis.ZV
	ZAdd(key string, members ...*redis.Z) int64
	ZAddOpt(key string, opts []string, members ...*redis.Z) (int64, error)
	ZCard(key string) int64
	ZRem(key string, members ...interface{}) int64
}

type IEs interface {
	GetClient() *es.Client
	Single(method, uri string, body []byte, timeout time.Duration) ([]byte, error)
	Batch(action, head, body string) error
}
