package adapter

import (
    "github.com/globalsign/mgo"
    "github.com/globalsign/mgo/bson"
    "github.com/pinguo/pgo2"
    "github.com/pinguo/pgo2/client/mongo"
    "github.com/pinguo/pgo2/iface"
)

func init() {
    container := pgo2.App().Container()
    container.Bind(&Mongo{})
}

// NewMongo of Mongo Client, add context support.
// usage: mongo := this.GetObject(Mongo.New(db, coll)).(adapter.IMongo)/(*adapter.Mongo)
func NewMongo(db, coll string, componentId ...string) *Mongo{
    id := DefaultMongoId
    if len(componentId) > 0 {
        id = componentId[0]
    }

    m := &Mongo{}

    m.client = pgo2.App().Component(id, mongo.New).(*mongo.Client)
    m.db = db
    m.coll = coll

    return m
    
}

// NewMongoPool of Mongo Client from pool, add context support.
// usage: mongo := this.GetObjectPool(Mongo.New,db, coll)).(adapter.IMongo)/(*adapter.Mongo)
func NewMongoPool(ctr iface.IContext, args ...interface{}) iface.IObject{
    if len(args) <2 {
        panic("need db and coll")
    }

    id := DefaultMongoId
    if len(args) > 2 {
        id = args[2].(string)
        if id == "" {
            panic("id must string")
        }
    }

    db := args[0].(string)
    coll := args[1].(string)

    if db == "" || coll == "" {
        panic("db and coll must string")
    }

    m := pgo2.App().GetObjPool(MongoClass, ctr).(*Mongo)

    m.client = pgo2.App().Component(id, mongo.New).(*mongo.Client)
    m.db = db
    m.coll = coll

    return m

}

type Mongo struct {
    pgo2.Object
    client *mongo.Client
    db     string
    coll   string
}



func (m *Mongo) GetClient() *mongo.Client {
    return m.client
}

// FindOne retrieve the first document that match the query,
// query can be a map or bson compatible struct, such as bson.M or properly typed map,
// nil query is equivalent to empty query such as bson.M{}.
// result is pointer to interface{}, map, bson.M or bson compatible struct, if interface{} type
// is provided, the output result is a bson.M.
// options provided optional query option listed as follows:
// fields: bson.M, set output fields, eg. bson.M{"_id":0, "name":1},
// sort: string or []string, set sort order, eg. "key1" or []string{"key1", "-key2"},
// skip: int, set skip number, eg. 100,
// limit: int, set result limit, eg. 1,
// hint: string or []string, set index hint, eg. []string{"key1", "key2"}
//
// for example:
//      var v1 interface{} // type of output v1 is bson.M
//      m.FindOne(bson.M{"_id":"k1"}, &v1)
//
//      var v2 struct {
//          Id    string `bson:"_id"`
//          Name  string `bson:"name"`
//          Value string `bson:"value"`
//      }
//      m.FindOne(bson.M{"_id": "k1"}, &v2)
func (m *Mongo) FindOne(query interface{}, result interface{}, options ...bson.M) error {
    profile := "Mongo.FindOne"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStop(profile)

    session := m.client.GetSession()
    defer session.Close()

    q := session.DB(m.db).C(m.coll).Find(query)
    m.applyQueryOptions(q, options)

    e := q.One(result)
    if e == nil || e == mgo.ErrNotFound {
        return nil
    }

    return e
}

// FindAll retrieve all documents that match the query,
// param result must be a slice(interface{}, map, bson.M or bson compatible struct)
// other params see FindOne()
func (m *Mongo) FindAll(query interface{}, result interface{}, options ...bson.M) error {
    profile := "Mongo.FindAll"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStop(profile)

    session := m.client.GetSession()
    defer session.Close()

    q := session.DB(m.db).C(m.coll).Find(query)
    m.applyQueryOptions(q, options)

    e := q.All(result)
    if e == nil || e == mgo.ErrNotFound {
        return nil
    }

    return e
}

// FindAndModify execute findAndModify command, which allows atomically update or remove one document,
// param change specify the change operation, eg. mgo.Change{Update:bson.M{"$inc": bson.M{"n":1}}, ReturnNew:true},
// other params see FindOne()
func (m *Mongo) FindAndModify(query interface{}, change mgo.Change, result interface{}, options ...bson.M) error {
    profile := "Mongo.FindAndModify"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStop(profile)

    session := m.client.GetSession()
    defer session.Close()

    q := session.DB(m.db).C(m.coll).Find(query)
    m.applyQueryOptions(q, options)

    _, e := q.Apply(change, result)
    if e == nil || e == mgo.ErrNotFound {
        return nil
    }

    return e
}

// FindDistinct retrieve distinct values for the param key,
// param result must be a slice,
// other params see FindOne()
func (m *Mongo) FindDistinct(query interface{}, key string, result interface{}, options ...bson.M) error {
    profile := "Mongo.FindDistinct"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStop(profile)

    session := m.client.GetSession()
    defer session.Close()

    q := session.DB(m.db).C(m.coll).Find(query)
    m.applyQueryOptions(q, options)

    e := q.Distinct(key, result)
    if e == nil || e == mgo.ErrNotFound {
        return nil
    }

    return e
}

// InsertOne insert one document into collection,
// param doc can be a map, bson.M, bson compatible struct,
// for example:
//      m.InsertOne(bson.M{"field1":"value1", "field2":"value2"})
//
//      doc := struct {
//          Field1 string `bson:"field1"`
//          Field2 string `bson:"field2"`
//      } {"value1", "value2"}
//      m.InsertOne(doc)
func (m *Mongo) InsertOne(doc interface{}) error {
    profile := "Mongo.InsertOne"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStop(profile)

    session := m.client.GetSession()
    defer session.Close()

    e := session.DB(m.db).C(m.coll).Insert(doc)

    return e
}

// InsertAll insert all documents provided by params docs into collection,
// for example:
//      docs := []interface{}{
//          bson.M{"_id":1, "name":"v1"},
//          bson.M{"_id":2, "name":"v2"},
//      }
//      m.InsertAll(docs)
func (m *Mongo) InsertAll(docs []interface{}) error {
    profile := "Mongo.InsertAll"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStop(profile)

    session := m.client.GetSession()
    defer session.Close()

    e := session.DB(m.db).C(m.coll).Insert(docs...)

    return e
}

// UpdateOne update one document that match the query,
// mgo.ErrNotFound is returned if a document not found,
// a value of *LastError is returned if other error occurred.
func (m *Mongo) UpdateOne(query interface{}, update interface{}) error {
    profile := "Mongo.UpdateOne"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStop(profile)

    session := m.client.GetSession()
    defer session.Close()

    e := session.DB(m.db).C(m.coll).Update(query, update)
    if e != nil && e != mgo.ErrNotFound {
        m.Context().Error(profile + " error, " + e.Error())
    }

    return e
}

// UpdateAll update all documents that match the query,
// see UpdateOne()
func (m *Mongo) UpdateAll(query interface{}, update interface{}) error {
    profile := "Mongo.UpdateAll"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStop(profile)

    session := m.client.GetSession()
    defer session.Close()

    _, e := session.DB(m.db).C(m.coll).UpdateAll(query, update)
    if e != nil && e != mgo.ErrNotFound {
        m.Context().Error(profile + " error, " + e.Error())
    }

    return e
}

// UpdateOrInsert update a existing document that match the query,
// or insert a new document base on the update document if no document match,
// an error of *LastError is returned if error is detected.
func (m *Mongo) UpdateOrInsert(query interface{}, update interface{}) error {
    profile := "Mongo.UpdateOrInsert"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStop(profile)

    session := m.client.GetSession()
    defer session.Close()

    _, e := session.DB(m.db).C(m.coll).Upsert(query, update)
    if e != nil {
        m.Context().Error(profile + " error, " + e.Error())
    }

    return e
}

// DeleteOne delete one document that match the query.
func (m *Mongo) DeleteOne(query interface{}) error {
    profile := "Mongo.DeleteOne"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStop(profile)

    session := m.client.GetSession()
    defer session.Close()

    e := session.DB(m.db).C(m.coll).Remove(query)
    if e != nil && e != mgo.ErrNotFound {
        m.Context().Error(profile + " error, " + e.Error())
    }

    return e
}

// DeleteAll delete all documents that match the query.
func (m *Mongo) DeleteAll(query interface{}) error {
    profile := "Mongo.DeleteAll"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStop(profile)

    session := m.client.GetSession()
    defer session.Close()

    _, e := session.DB(m.db).C(m.coll).RemoveAll(query)
    if e != nil && e != mgo.ErrNotFound {
        m.Context().Error(profile + " error, " + e.Error())
    }

    return e
}

// Count return the count of documents match the query.
func (m *Mongo) Count(query interface{}, options ...bson.M) (int, error) {
    profile := "Mongo.Count"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStop(profile)

    session := m.client.GetSession()
    defer session.Close()

    q := session.DB(m.db).C(m.coll).Find(query)
    m.applyQueryOptions(q, options)

    n, e := q.Count()
    if e != nil {
        m.Context().Error(profile + " error, " + e.Error())
    }

    return n, e
}

// PipeOne execute aggregation queries and get the first item from result set.
// param pipeline must be a slice, such as []bson.M,
// param result is a pointer to interface{}, map, bson.M or bson compatible struct.
// for example:
//      pipeline := []bson.M{
//          bson.M{"$match": bson.M{"status":"A"}},
//          bson.M{"$group": bson.M{"_id":"$field1", "total":"$field2"}},
//      }
//      m.PipeOne(pipeline, &result)
func (m *Mongo) PipeOne(pipeline interface{}, result interface{}) error {
    profile := "Mongo.PipeOne"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStop(profile)

    session := m.client.GetSession()
    defer session.Close()

    e := session.DB(m.db).C(m.coll).Pipe(pipeline).One(result)
    if e != nil && e != mgo.ErrNotFound {
        m.Context().Error(profile + " error, " + e.Error())
    }

    return e
}

// PipeAll execute aggregation queries and get all item from result set.
// param result must be slice(interface{}, map, bson.M or bson compatible struct).
// see PipeOne().
func (m *Mongo) PipeAll(pipeline interface{}, result interface{}) error {
    profile := "Mongo.PipeAll"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStop(profile)

    session := m.client.GetSession()
    defer session.Close()

    e := session.DB(m.db).C(m.coll).Pipe(pipeline).All(result)
    if e != nil && e != mgo.ErrNotFound {
        m.Context().Error(profile + " error, " + e.Error())
    }

    return e
}

// MapReduce execute map/reduce job that match the query.
// param result is a slice(interface{}, map, bson.M, bson compatible struct),
// param query and options see FindOne().
// for example:
//      job := &mgo.MapReduce{
//          Map: "function() { emit(this.n, 1) }",
//          Reduce: "function(key, values) { return Array.sum(values) }",
//      }
//      result := []bson.M{}
//      m.MapReduce(query, job, &result)
func (m *Mongo) MapReduce(query interface{}, job *mgo.MapReduce, result interface{}, options ...bson.M) error {
    profile := "Mongo.MapReduce"
    m.Context().ProfileStart(profile)
    defer m.Context().ProfileStop(profile)

    session := m.client.GetSession()
    defer session.Close()

    q := session.DB(m.db).C(m.coll).Find(query)
    m.applyQueryOptions(q, options)

    _, e := q.MapReduce(job, result)
    if e != nil && e != mgo.ErrNotFound {
        m.Context().Error(profile + " error, " + e.Error())
    }

    return e
}

func (m *Mongo) applyQueryOptions(q *mgo.Query, options []bson.M) {
    if len(options) == 0 {
        return
    }

    for key, opt := range options[0] {
        switch key {
        case "fields":
            if fields, ok := opt.(bson.M); ok {
                q.Select(fields)
            }

        case "sort":
            if arr, ok := opt.([]string); ok {
                q.Sort(arr...)
            } else if str, ok := opt.(string); ok {
                q.Sort(str)
            }

        case "skip":
            if skip, ok := opt.(int); ok {
                q.Skip(skip)
            }

        case "limit":
            if limit, ok := opt.(int); ok {
                q.Limit(limit)
            }

        case "hint":
            if arr, ok := opt.([]string); ok {
                q.Hint(arr...)
            } else if str, ok := opt.(string); ok {
                q.Hint(str)
            }

        default:
            panic(ErrInvalidOpt + key)
        }
    }
}
