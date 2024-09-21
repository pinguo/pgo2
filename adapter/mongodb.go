package adapter

import (
	"context"
	"fmt"

	"github.com/globalsign/mgo/bson"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/pinguo/pgo2"
	"github.com/pinguo/pgo2/client/mongodb"

	opts "github.com/qiniu/qmgo/options"
)

var MongodbClass string
var MongodbCQueryClass string
var MongodbCAggregate string

func init() {
	container := pgo2.App().Container()
	MongodbClass = container.Bind(&Mongodb{})
	MongodbCQueryClass = container.Bind(&CQuery{})
	MongodbCAggregate = container.Bind(&CAggregate{})
}

// NewMongodb of Mongodb Client, add context support.
// usage: Mongodb := this.GetObj(adapter.NewMongodb(db, coll)).(adapter.IMongodb)/(*adapter.Mongodb)
func NewMongodb(db, coll string, componentId ...string) *Mongodb {
	id := DefaultMongodbId
	if len(componentId) > 0 {
		id = componentId[0]
	}

	m := &Mongodb{}

	m.client = pgo2.App().Component(id, mongodb.New).(*mongodb.Client)
	m.db = db
	m.coll = coll

	return m

}

type Mongodb struct {
	pgo2.Object
	client *mongodb.Client
	*qmgo.Collection
	db   string
	coll string
}

// GetObjBox fetch is performed automatically
func (m *Mongodb) Prepare(db, coll string, componentId ...string) {
	if db == "" || coll == "" {
		panic("db and coll can not empty")
	}

	id := DefaultMongodbId
	if len(componentId) > 0 {
		id = componentId[0]
		if id == "" {
			panic("id must string")
		}
	}

	m.client = pgo2.App().Component(id, mongodb.New).(*mongodb.Client)
	m.db = db
	m.coll = coll

}

func (m *Mongodb) GetClient() *mongodb.Client {
	return m.client
}

// InsertOne insert one document into the collection
// If InsertHook in opts is set, hook works on it, otherwise hook try the doc as hook
// Reference: https://docs.mongodb.com/manual/reference/command/insert/
func (m *Mongodb) InsertOne(doc interface{}, opts ...opts.InsertOneOptions) (result *qmgo.InsertOneResult, err error) {
	profile := "mongodb.InsertOne"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)
	ctx, cFunc := context.WithTimeout(context.Background(), m.client.WriteTimeout())
	defer cFunc()
	// ctx := context.Background()
	return m.client.MClient.Database(m.db).Collection(m.coll).InsertOne(ctx, doc, opts...)
}

func (m *Mongodb) InsertOneCtx(ctx context.Context, doc interface{}, opts ...opts.InsertOneOptions) (result *qmgo.InsertOneResult, err error) {
	profile := "mongodb.InsertOneCtx"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)

	return m.client.MClient.Database(m.db).Collection(m.coll).InsertOne(ctx, doc, opts...)
}

// InsertMany executes an insert command to insert multiple documents into the collection.
// If InsertHook in opts is set, hook works on it, otherwise hook try the doc as hook
// Reference: https://docs.mongodb.com/manual/reference/command/insert/
func (m *Mongodb) InsertMany(docs interface{}, opts ...opts.InsertManyOptions) (result *qmgo.InsertManyResult, err error) {
	profile := "mongodb.InsertMany"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)
	ctx, cFunc := context.WithTimeout(context.Background(), m.client.WriteTimeout())
	defer cFunc()

	return m.client.MClient.Database(m.db).Collection(m.coll).InsertMany(ctx, docs, opts...)
}

func (m *Mongodb) InsertManyCtx(ctx context.Context, docs interface{}, opts ...opts.InsertManyOptions) (result *qmgo.InsertManyResult, err error) {
	profile := "mongodb.InsertManyCtx"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)

	return m.client.MClient.Database(m.db).Collection(m.coll).InsertMany(ctx, docs, opts...)
}

// Upsert updates one documents if filter match, inserts one document if filter is not match, Error when the filter is invalid
// The replacement parameter must be a document that will be used to replace the selected document. It cannot be nil
// and cannot contain any update operators
// Reference: https://docs.mongodb.com/manual/reference/operator/update/
// If replacement has "_id" field and the document is exist, please initial it with existing id(even with Qmgo default field feature).
// Otherwise "the (immutable) field '_id' altered" error happens.
func (m *Mongodb) Upsert(filter interface{}, replacement interface{}, opts ...opts.UpsertOptions) (result *qmgo.UpdateResult, err error) {
	profile := "mongodb.Upsert"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)
	ctx, cFunc := context.WithTimeout(context.Background(), m.client.WriteTimeout())
	defer cFunc()

	return m.client.MClient.Database(m.db).Collection(m.coll).Upsert(ctx, filter, replacement, opts...)
}

func (m *Mongodb) UpsertCtx(ctx context.Context, filter interface{}, replacement interface{}, opts ...opts.UpsertOptions) (result *qmgo.UpdateResult, err error) {
	profile := "mongodb.UpsertCtx"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)

	return m.client.MClient.Database(m.db).Collection(m.coll).Upsert(ctx, filter, replacement, opts...)
}

// UpsertId updates one documents if id match, inserts one document if id is not match and the id will inject into the document
// The replacement parameter must be a document that will be used to replace the selected document. It cannot be nil
// and cannot contain any update operators
// Reference: https://docs.mongodb.com/manual/reference/operator/update/
func (m *Mongodb) UpsertId(id interface{}, replacement interface{}, opts ...opts.UpsertOptions) (result *qmgo.UpdateResult, err error) {
	profile := "mongodb.UpsertId"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)
	ctx, cFunc := context.WithTimeout(context.Background(), m.client.WriteTimeout())
	defer cFunc()

	return m.client.MClient.Database(m.db).Collection(m.coll).UpsertId(ctx, id, replacement, opts...)
}

func (m *Mongodb) UpsertIdCtx(ctx context.Context, id interface{}, replacement interface{}, opts ...opts.UpsertOptions) (result *qmgo.UpdateResult, err error) {
	profile := "mongodb.UpsertIdCtx"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)

	return m.client.MClient.Database(m.db).Collection(m.coll).Upsert(ctx, id, replacement, opts...)
}

func handleOptions(query qmgo.QueryI, options ...bson.M) (qmgo.QueryI, error) {
	opts := make(map[string]interface{})
	for _, opt := range options {
		if opt == nil {
			continue
		}
		for k, v := range opt {
			opts[k] = v
		}
	}
	if sort, ok := opts["sort"]; ok {
		switch sort.(type) {
		case string:
			query.Sort(sort.(string))
		case []string:
			query.Sort(sort.([]string)...)
		default:
			return nil, fmt.Errorf("invalid mongo sort:%#v", sort)
		}
	}
	number := func(name string) (int64, error) {
		rl := int64(0)
		v, ok := opts[name]
		if !ok {
			return 0, nil
		}
		switch v.(type) {
		case int:
			rl = int64(v.(int))
		case int32:
			rl = int64(v.(int32))
		case int64:
			rl = v.(int64)
		default:
			return 0, fmt.Errorf("invalid mongo limit: %#v", v)
		}
		return rl, nil
	}
	if limit, err := number("limit"); err != nil {
		return nil, err
	} else if limit > 0 {
		query.Limit(limit)
	}

	if skip, err := number("skip"); err != nil {
		return nil, err
	} else if skip > 0 {
		query.Skip(skip)
	}
	return query, nil
}

// DeleteAll delete all doc
func (m *Mongodb) DeleteAll(filter interface{}) error {
	profile := "mongodb.DeleteAll"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)
	ctx, cFunc := context.WithTimeout(context.Background(), m.client.WriteTimeout())
	defer cFunc()
	col := m.client.MClient.Database(m.db).Collection(m.coll)
	_, err := col.RemoveAll(ctx, filter)
	return err
}

// Count count docs
func (m *Mongodb) Count(filter interface{}) (int, error) {
	profile := "mongodb.Count"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)
	ctx, cFunc := context.WithTimeout(context.Background(), m.client.WriteTimeout())
	defer cFunc()
	col := m.client.MClient.Database(m.db).Collection(m.coll)
	n, err := col.Find(ctx, filter).Count()
	return int(n), err
}

// DeleteOne delete one, 暂时没有one
func (m *Mongodb) DeleteOne(filter interface{}) error {
	profile := "mongodb.DeleteOne"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)
	ctx, cFunc := context.WithTimeout(context.Background(), m.client.WriteTimeout())
	defer cFunc()
	col := m.client.MClient.Database(m.db).Collection(m.coll)
	return col.Remove(ctx, filter, opts.RemoveOptions{
		DeleteOptions: &options.DeleteOptions{},
	})
}

func (m *Mongodb) InsertAll(docs interface{}) error {
	profile := "mongodb.InsertAll"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)
	ctx, cFunc := context.WithTimeout(context.Background(), m.client.WriteTimeout())
	defer cFunc()
	col := m.client.MClient.Database(m.db).Collection(m.coll)
	_, err := col.InsertMany(ctx, docs)
	return err
}

func (m *Mongodb) UpdateOrInsert(query interface{}, doc interface{}) error {
	profile := "mongodb.UpdateOrInsert"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)
	ctx, cFunc := context.WithTimeout(context.Background(), m.client.WriteTimeout())
	defer cFunc()
	col := m.client.MClient.Database(m.db).Collection(m.coll)
	_, err := col.Upsert(ctx, query, doc)
	return err
}

func (m *Mongodb) PipeAll(query interface{}, doc interface{}) error {
	profile := "mongodb.PipeAll"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)
	ctx, cFunc := context.WithTimeout(context.Background(), m.client.WriteTimeout())
	defer cFunc()
	col := m.client.MClient.Database(m.db).Collection(m.coll)
	ag := col.Aggregate(ctx, query)
	return ag.All(doc)
}

// FindAll query all doc
func (m *Mongodb) FindAll(filter interface{}, doc interface{}, options ...bson.M) (err error) {
	profile := "mongodb.FindOne"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)
	ctx, cFunc := context.WithTimeout(context.Background(), m.client.WriteTimeout())
	defer cFunc()
	col := m.client.MClient.Database(m.db).Collection(m.coll)
	findOpts := opts.FindOptions{}
	cursor, err := handleOptions(col.Find(ctx, filter, findOpts), options...)
	if err != nil {
		return err
	}
	return cursor.All(doc)
}

// FindOne query one doc
func (m *Mongodb) FindOne(filter interface{}, doc interface{}, options ...bson.M) (err error) {
	profile := "mongodb.FindOne"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)
	ctx, cFunc := context.WithTimeout(context.Background(), m.client.WriteTimeout())
	defer cFunc()
	col := m.client.MClient.Database(m.db).Collection(m.coll)
	findOpts := opts.FindOptions{}
	cursor, err := handleOptions(col.Find(ctx, filter, findOpts), options...)
	if err != nil {
		return err
	}
	return cursor.One(doc)
}

// UpdateOne executes an update command to update at most one document in the collection.
// Reference: https://docs.mongodb.com/manual/reference/operator/update/
func (m *Mongodb) UpdateOne(filter interface{}, update interface{}, opts ...opts.UpdateOptions) (err error) {
	profile := "mongodb.UpdateOne"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)
	ctx, cFunc := context.WithTimeout(context.Background(), m.client.WriteTimeout())
	defer cFunc()

	return m.client.MClient.Database(m.db).Collection(m.coll).UpdateOne(ctx, filter, update, opts...)
}

func (m *Mongodb) UpdateOneCtx(ctx context.Context, filter interface{}, update interface{}, opts ...opts.UpdateOptions) (err error) {
	profile := "mongodb.UpdateOneCtx"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)

	return m.client.MClient.Database(m.db).Collection(m.coll).UpdateOne(ctx, filter, update, opts...)
}

// UpdateId executes an update command to update at most one document in the collection.
// Reference: https://docs.mongodb.com/manual/reference/operator/update/
func (m *Mongodb) UpdateId(id interface{}, update interface{}, opts ...opts.UpdateOptions) (err error) {
	profile := "mongodb.UpdateId"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)
	ctx, cFunc := context.WithTimeout(context.Background(), m.client.WriteTimeout())
	defer cFunc()

	return m.client.MClient.Database(m.db).Collection(m.coll).UpdateId(ctx, id, update, opts...)
}

func (m *Mongodb) UpdateIdCtx(ctx context.Context, id interface{}, update interface{}, opts ...opts.UpdateOptions) (err error) {
	profile := "mongodb.UpdateIdCtx"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)

	return m.client.MClient.Database(m.db).Collection(m.coll).UpdateId(ctx, id, update, opts...)
}

// UpdateAll executes an update command to update documents in the collection.
// The matchedCount is 0 in UpdateResult if no document updated
// Reference: https://docs.mongodb.com/manual/reference/operator/update/
func (m *Mongodb) UpdateAll(filter interface{}, update interface{}, opts ...opts.UpdateOptions) (result *qmgo.UpdateResult, err error) {
	profile := "mongodb.UpdateAll"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)
	ctx, cFunc := context.WithTimeout(context.Background(), m.client.WriteTimeout())
	defer cFunc()

	return m.client.MClient.Database(m.db).Collection(m.coll).UpdateAll(ctx, filter, update, opts...)
}

func (m *Mongodb) UpdateAllCtx(ctx context.Context, filter interface{}, update interface{}, opts ...opts.UpdateOptions) (result *qmgo.UpdateResult, err error) {
	profile := "mongodb.UpdateAllCtx"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)

	return m.client.MClient.Database(m.db).Collection(m.coll).UpdateAll(ctx, filter, update, opts...)
}

// ReplaceOne executes an update command to update at most one document in the collection.
// If UpdateHook in opts is set, hook works on it, otherwise hook try the doc as hook
// Expect type of the doc is the define of user's document
func (m *Mongodb) ReplaceOne(filter interface{}, doc interface{}, opts ...opts.ReplaceOptions) (err error) {
	profile := "mongodb.ReplaceOne"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)
	ctx, cFunc := context.WithTimeout(context.Background(), m.client.WriteTimeout())
	defer cFunc()

	return m.client.MClient.Database(m.db).Collection(m.coll).ReplaceOne(ctx, filter, doc, opts...)
}

func (m *Mongodb) ReplaceOneCtx(ctx context.Context, filter interface{}, doc interface{}, opts ...opts.ReplaceOptions) (err error) {
	profile := "mongodb.ReplaceOneCtx"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)

	return m.client.MClient.Database(m.db).Collection(m.coll).ReplaceOne(ctx, filter, doc, opts...)
}

// Remove executes a delete command to delete at most one document from the collection.
// if filter is bson.M{}，DeleteOne will delete one document in collection
// Reference: https://docs.mongodb.com/manual/reference/command/delete/
func (m *Mongodb) Remove(filter interface{}, opts ...opts.RemoveOptions) (err error) {
	profile := "mongodb.Remove"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)
	ctx, cFunc := context.WithTimeout(context.Background(), m.client.WriteTimeout())
	defer cFunc()

	return m.client.MClient.Database(m.db).Collection(m.coll).Remove(ctx, filter, opts...)
}

func (m *Mongodb) RemoveCtx(ctx context.Context, filter interface{}, opts ...opts.RemoveOptions) (err error) {
	profile := "mongodb.RemoveCtx"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)

	return m.client.MClient.Database(m.db).Collection(m.coll).Remove(ctx, filter, opts...)
}

// RemoveId executes a delete command to delete at most one document from the collection.
func (m *Mongodb) RemoveId(id interface{}, opts ...opts.RemoveOptions) (err error) {
	profile := "mongodb.RemoveId"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)
	ctx, cFunc := context.WithTimeout(context.Background(), m.client.WriteTimeout())
	defer cFunc()

	return m.client.MClient.Database(m.db).Collection(m.coll).Remove(ctx, id, opts...)
}

func (m *Mongodb) RemoveIdCtx(ctx context.Context, id interface{}, opts ...opts.RemoveOptions) (err error) {
	profile := "mongodb.RemoveIdCtx"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)

	return m.client.MClient.Database(m.db).Collection(m.coll).Remove(ctx, id, opts...)
}

func (m *Mongodb) RemoveAll(filter interface{}, opts ...opts.RemoveOptions) (result *qmgo.DeleteResult, err error) {
	profile := "mongodb.RemoveAll"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)
	ctx, cFunc := context.WithTimeout(context.Background(), m.client.WriteTimeout())
	defer cFunc()

	return m.client.MClient.Database(m.db).Collection(m.coll).RemoveAll(ctx, filter, opts...)
}

func (m *Mongodb) RemoveAllCtx(ctx context.Context, filter interface{}, opts ...opts.RemoveOptions) (result *qmgo.DeleteResult, err error) {
	profile := "mongodb.RemoveAllCtx"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)

	return m.client.MClient.Database(m.db).Collection(m.coll).RemoveAll(ctx, filter, opts...)
}

// Aggregate executes an aggregate command against the collection and returns a AggregateI to get resulting documents.
func (m *Mongodb) Aggregate(pipeline interface{}) IMongodbAggregate {
	ctx, cFunc := context.WithTimeout(context.Background(), m.client.ReadTimeout())
	aggregate := m.client.MClient.Database(m.db).Collection(m.coll).Aggregate(ctx, pipeline)

	return m.GetObjBoxCtx(m.Context(), MongodbCAggregate, aggregate, cFunc).(IMongodbAggregate)

}

func (m *Mongodb) AggregateCtx(ctx context.Context, pipeline interface{}) IMongodbAggregate {
	aggregate := m.client.MClient.Database(m.db).Collection(m.coll).Aggregate(ctx, pipeline)

	return m.GetObjBoxCtx(m.Context(), MongodbCAggregate, aggregate).(IMongodbAggregate)
}

// EnsureIndexes Deprecated
// Recommend to use CreateIndexes / CreateOneIndex for more function)
// EnsureIndexes creates unique and non-unique indexes in collection
// the combination of indexes is different from CreateIndexes:
// if uniques/indexes is []string{"name"}, means create index "name"
// if uniques/indexes is []string{"name,-age","uid"} means create Compound indexes: name and -age, then create one index: uid
func (m *Mongodb) EnsureIndexes(uniques []string, indexes []string) (err error) {
	profile := "mongodb.EnsureIndexes"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)
	ctx, cFunc := context.WithTimeout(context.Background(), m.client.WriteTimeout())
	defer cFunc()

	return m.client.MClient.Database(m.db).Collection(m.coll).EnsureIndexes(ctx, uniques, indexes)
}

func (m *Mongodb) EnsureIndexesCtx(ctx context.Context, uniques []string, indexes []string) (err error) {
	profile := "mongodb.EnsureIndexesCtx"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)

	return m.client.MClient.Database(m.db).Collection(m.coll).EnsureIndexes(ctx, uniques, indexes)
}

// CreateIndexes creates multiple indexes in collection
// If the Key in opts.IndexModel is []string{"name"}, means create index: name
// If the Key in opts.IndexModel is []string{"name","-age"} means create Compound indexes: name and -age
func (m *Mongodb) CreateIndexes(indexes []opts.IndexModel) (err error) {
	profile := "mongodb.CreateIndexes"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)
	ctx, cFunc := context.WithTimeout(context.Background(), m.client.WriteTimeout())
	defer cFunc()

	return m.client.MClient.Database(m.db).Collection(m.coll).CreateIndexes(ctx, indexes)
}

func (m *Mongodb) CreateIndexesCtx(ctx context.Context, indexes []opts.IndexModel) (err error) {
	profile := "mongodb.CreateIndexesCtx"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)

	return m.client.MClient.Database(m.db).Collection(m.coll).CreateIndexes(ctx, indexes)
}

// CreateOneIndex creates one index
// If the Key in opts.IndexModel is []string{"name"}, means create index name
// If the Key in opts.IndexModel is []string{"name","-age"} means create Compound index: name and -age
func (m *Mongodb) CreateOneIndex(index opts.IndexModel) error {
	profile := "mongodb.CreateOneIndex"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)
	ctx, cFunc := context.WithTimeout(context.Background(), m.client.WriteTimeout())
	defer cFunc()

	return m.client.MClient.Database(m.db).Collection(m.coll).CreateOneIndex(ctx, index)
}

func (m *Mongodb) CreateOneIndexCtx(ctx context.Context, index opts.IndexModel) (err error) {
	profile := "mongodb.CreateOneIndexCtx"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)

	return m.client.MClient.Database(m.db).Collection(m.coll).CreateOneIndex(ctx, index)
}

// DropAllIndexes drop all indexes on the collection except the index on the _id field
// if there is only _id field index on the collection, the function call will report an error
func (m *Mongodb) DropAllIndexes() (err error) {
	profile := "mongodb.DropAllIndexes"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)
	ctx, cFunc := context.WithTimeout(context.Background(), m.client.WriteTimeout())
	defer cFunc()

	return m.client.MClient.Database(m.db).Collection(m.coll).DropAllIndexes(ctx)
}

func (m *Mongodb) DropAllIndexesCtx(ctx context.Context) (err error) {
	profile := "mongodb.DropAllIndexesCtx"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)

	return m.client.MClient.Database(m.db).Collection(m.coll).DropAllIndexes(ctx)
}

// DropIndex drop indexes in collection, indexes that be dropped should be in line with inputting indexes
// The indexes is []string{"name"} means drop index: name
// The indexes is []string{"name","-age"} means drop Compound indexes: name and -age
func (m *Mongodb) DropIndex(indexes []string) error {
	profile := "mongodb.DropIndex"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)
	ctx, cFunc := context.WithTimeout(context.Background(), m.client.WriteTimeout())
	defer cFunc()

	return m.client.MClient.Database(m.db).Collection(m.coll).DropIndex(ctx, indexes)
}

func (m *Mongodb) DropIndexCtx(ctx context.Context, indexes []string) error {
	profile := "mongodb.DropIndexCtx"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)

	return m.client.MClient.Database(m.db).Collection(m.coll).DropIndex(ctx, indexes)
}

// DropCollection drops collection
// it's safe even collection is not exists
func (m *Mongodb) DropCollection() error {
	profile := "mongodb.DropCollection"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)
	ctx, cFunc := context.WithTimeout(context.Background(), m.client.WriteTimeout())
	defer cFunc()

	return m.client.MClient.Database(m.db).Collection(m.coll).DropCollection(ctx)
}

func (m *Mongodb) DropCollectionCtx(ctx context.Context) error {
	profile := "mongodb.DropCollectionCtx"
	m.Context().ProfileStart(profile)
	defer m.Context().ProfileStop(profile)

	return m.client.MClient.Database(m.db).Collection(m.coll).DropCollection(ctx)
}

func (m *Mongodb) CloneCollection() (*mongo.Collection, error) {

	return m.client.MClient.Database(m.db).Collection(m.coll).CloneCollection()
}

func (m *Mongodb) GetCollectionName() string {

	return m.client.MClient.Database(m.db).Collection(m.coll).GetCollectionName()
}

func (m *Mongodb) Session() (*qmgo.Session, error) {

	return m.client.MClient.Session()
}

func (m *Mongodb) DoTransaction(ctx context.Context, callback func(sessCtx context.Context) (interface{}, error)) (interface{}, error) {

	return m.client.MClient.DoTransaction(ctx, callback)
}

func (m *Mongodb) Find(filter interface{}, options ...opts.FindOptions) IMongodbQuery {
	ctx, cFunc := context.WithTimeout(context.Background(), m.client.ReadTimeout())
	//	defer cFunc()
	query := m.client.MClient.Database(m.db).Collection(m.coll).Find(ctx, filter, options...)

	return m.GetObjBoxCtx(m.Context(), MongodbCQueryClass, query, cFunc).(*CQuery)
}

func (m *Mongodb) FindCtx(ctx context.Context, filter interface{}, options ...opts.FindOptions) IMongodbQuery {
	query := m.client.MClient.Database(m.db).Collection(m.coll).Find(ctx, filter, options...)
	return m.GetObjBoxCtx(m.Context(), MongodbCQueryClass, query).(*CQuery)
}

type CQuery struct {
	qmgo.QueryI
	pgo2.Object
	ctxCancel context.CancelFunc
}

func (c *CQuery) Prepare(q qmgo.QueryI, dftCtxCancel ...context.CancelFunc) {
	c.QueryI = q
	if len(dftCtxCancel) > 0 {
		c.ctxCancel = dftCtxCancel[0]
	}
}

func (c *CQuery) Sort(fields ...string) qmgo.QueryI {
	q := c.QueryI.Sort(fields...)
	return c.GetObjBoxCtx(c.Context(), MongodbCQueryClass, q, c.ctxCancel).(*CQuery)
}

func (c *CQuery) Select(selector interface{}) qmgo.QueryI {
	q := c.QueryI.Select(selector)
	return c.GetObjBoxCtx(c.Context(), MongodbCQueryClass, q, c.ctxCancel).(*CQuery)
}

func (c *CQuery) Skip(n int64) qmgo.QueryI {
	q := c.QueryI.Skip(n)
	return c.GetObjBoxCtx(c.Context(), MongodbCQueryClass, q, c.ctxCancel).(*CQuery)
}

func (c *CQuery) Limit(n int64) qmgo.QueryI {
	q := c.QueryI.Limit(n)
	return c.GetObjBoxCtx(c.Context(), MongodbCQueryClass, q, c.ctxCancel).(*CQuery)
}

func (c *CQuery) Hint(hint interface{}) qmgo.QueryI {
	q := c.QueryI.Hint(hint)
	return c.GetObjBoxCtx(c.Context(), MongodbCQueryClass, q, c.ctxCancel).(*CQuery)
}

func (c *CQuery) One(result interface{}) error {
	profile := "mongodb.One"
	c.Context().ProfileStart(profile)
	defer c.Context().ProfileStop(profile)
	if c.ctxCancel != nil {
		defer c.ctxCancel()
	}

	return c.QueryI.One(result)
}

func (c *CQuery) All(result interface{}) error {
	profile := "mongodb.All"
	c.Context().ProfileStart(profile)
	defer c.Context().ProfileStop(profile)
	if c.ctxCancel != nil {
		defer c.ctxCancel()
	}

	return c.QueryI.All(result)
}

func (c *CQuery) Count() (n int64, err error) {
	profile := "mongodb.Count"
	c.Context().ProfileStart(profile)
	defer c.Context().ProfileStop(profile)
	if c.ctxCancel != nil {
		defer c.ctxCancel()
	}

	return c.QueryI.Count()
}

func (c *CQuery) Distinct(key string, result interface{}) error {
	profile := "mongodb.Distinct"
	c.Context().ProfileStart(profile)
	defer c.Context().ProfileStop(profile)
	if c.ctxCancel != nil {
		defer c.ctxCancel()
	}

	return c.QueryI.Distinct(key, result)
}

func (c *CQuery) Apply(change qmgo.Change, result interface{}) error {
	profile := "mongodb.Apply"
	c.Context().ProfileStart(profile)
	defer c.Context().ProfileStop(profile)
	if c.ctxCancel != nil {
		defer c.ctxCancel()
	}

	return c.QueryI.Apply(change, result)
}

type CAggregate struct {
	pgo2.Object
	qmgo.AggregateI
	ctxCancel context.CancelFunc
}

func (c *CAggregate) Prepare(q qmgo.AggregateI, dftCtxCancel ...context.CancelFunc) {
	c.AggregateI = q
	if len(dftCtxCancel) > 0 {
		c.ctxCancel = dftCtxCancel[0]
	}

}

func (c *CAggregate) All(results interface{}) error {
	profile := "mongodb.aggregate.All"
	c.Context().ProfileStart(profile)
	defer c.Context().ProfileStop(profile)
	if c.ctxCancel != nil {
		defer c.ctxCancel()
	}

	return c.AggregateI.All(results)
}

func (c *CAggregate) One(results interface{}) error {
	profile := "mongodb.aggregate.One"
	c.Context().ProfileStart(profile)
	defer c.Context().ProfileStop(profile)
	if c.ctxCancel != nil {
		defer c.ctxCancel()
	}

	return c.AggregateI.One(results)
}
