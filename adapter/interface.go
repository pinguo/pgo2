package adapter

import (
    "context"
    "database/sql"
    "net/http"
    "time"

    "github.com/globalsign/mgo"
    "github.com/globalsign/mgo/bson"
    "github.com/pinguo/pgo2/client/maxmind"
    "github.com/pinguo/pgo2/client/memcache"
    "github.com/pinguo/pgo2/client/phttp"
    "github.com/pinguo/pgo2/client/rabbitmq"
    "github.com/pinguo/pgo2/iface"
    "github.com/pinguo/pgo2/value"
    "github.com/streadway/amqp"
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
    Do(cmd string, args ...interface{}) interface{}
}

type IRabbitMq interface {
    SetPanicRecover(v bool)
    ExchangeDeclare()
    Publish(opCode string, data interface{}, dftOpUid ...string) bool
    GetConsumeChannelBox(queueName string, opCodes []string) *rabbitmq.ChannelBox
    Consume(queueName string, opCodes []string, limit int, autoAck, noWait, exclusive bool) <-chan amqp.Delivery
    DecodeBody(d amqp.Delivery, ret interface{}) error
    DecodeHeaders(d amqp.Delivery) *rabbitmq.RabbitHeaders
}

type IDb interface {
    GetDb(master bool) *sql.DB
    Begin(opts ...*sql.TxOptions) ITx
    BeginContext(ctx context.Context, opts *sql.TxOptions) ITx
    QueryOne(query string, args ...interface{}) IRow
    QueryOneContext(ctx context.Context, query string, args ...interface{}) IRow
    Query(query string, args ...interface{}) *sql.Rows
    QueryContext(ctx context.Context, query string, args ...interface{}) *sql.Rows
    Exec(query string, args ...interface{}) sql.Result
    ExecContext(ctx context.Context, query string, args ...interface{}) sql.Result
    Prepare(query string) IStmt
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
    Prepare(query string) IStmt
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