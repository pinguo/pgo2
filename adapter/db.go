package adapter

import (
    "context"
    "database/sql"
    "strings"
    "sync"
    "time"

    "github.com/pinguo/pgo2"
    "github.com/pinguo/pgo2/client/db"
    "github.com/pinguo/pgo2/iface"
    "github.com/pinguo/pgo2/util"
)

var (
    stmtPool sync.Pool
    rowPool  sync.Pool
    txPool   sync.Pool
)

func init() {
    stmtPool.New = func() interface{} { return &Stmt{} }
    rowPool.New = func() interface{} { return &Row{} }
    txPool.New = func() interface{} { return &Tx{} }
    container := pgo2.App().Container()
    container.Bind(&Db{})
}

type Db struct {
    base
    db *sql.DB

    master bool
}

// NewDb of db Client, add context support.
// usage: db := this.GetObj(adapter.NewDb()).(adapter.IDb)/(*adapter.Db)
func NewDb(componentId ...string) *Db {
    id := DefaultDbId
    if len(componentId) > 0 {
        id = componentId[0]
    }
    d := &Db{}
    d.client = pgo2.App().Component(id, db.New).(*db.Client)

    return d
}

// NewDbPool of db Client from pool, add context support.
// usage: db := this.GetObjPool(adapter.NewDbPool).(adapter.IDb)/(*adapter.Db)
func NewDbPool(ctr iface.IContext, componentId ...interface{}) iface.IObject {
    id := DefaultDbId
    if len(componentId) > 0 {
        id = componentId[0].(string)
    }

    d := pgo2.App().GetObjPool(DbClass, ctr).(*Db)

    d.client = pgo2.App().Component(id, db.New).(*db.Client)

    return d
}

func (d *Db) SetMaster(v bool) {
    d.master = v
}

func (d *Db) Master() bool {
    return d.master
}

func (d *Db) GetClient() *db.Client {
    return d.client
}

func (d *Db) GetDb(master bool) *sql.DB {
    // reuse previous db instance for read
    if !master && d.db != nil {
        return d.db
    }
    d.db = d.client.GetDb(master)
    return d.db
}

// Begin start a transaction with default timeout context and optional opts,
// if opts is nil, default driver option will be used.
func (d *Db) Begin(opts ...*sql.TxOptions) ITx {
    opts = append(opts, nil)
    ctx, _ := context.WithTimeout(context.Background(), DefaultDbTimeout)
    return d.BeginContext(ctx, opts[0])
}

// BeginContext start a transaction with specified context and optional opts,
// if opts is nil, default driver option will be used.
func (d *Db) BeginContext(ctx context.Context, opts *sql.TxOptions) ITx {
    if tx, e := d.GetDb(true).BeginTx(ctx, opts); e != nil {
        d.Context().Error("Db.Begin error, " + e.Error())
        return nil
    } else {
        txWrapper := txPool.Get().(*Tx)
        txWrapper.SetContext(d.Context())
        txWrapper.init(tx, d.client)
        return txWrapper
    }
}

// QueryOne perform one row query using a default timeout context,
// and always returns a non-nil value, Errors are deferred until
// Row's Scan method is called.
func (d *Db) QueryOne(query string, args ...interface{}) IRow {
    ctx, _ := context.WithTimeout(context.Background(), DefaultDbTimeout)
    return d.QueryOneContext(ctx, query, args...)
}

// QueryOneContext perform one row query using a specified context,
// and always returns a non-nil value, Errors are deferred until
// Row's Scan method is called.
func (d *Db) QueryOneContext(ctx context.Context, query string, args ...interface{}) IRow {
    start := time.Now()
    defer func() {
        elapse := time.Since(start)
        d.profileAdd(nil, "Db.QueryOne", elapse, query, args...)

        if elapse >= d.client.SlowLogTime() && d.client.SlowLogTime() > 0 {
            d.Context().Warn("Db.QueryOne slow, elapse:%s, query:%s, args:%v", elapse, query, args)
        }
    }()

    row := d.GetDb(d.master).QueryRowContext(ctx, query, args...)

    // wrap row for profile purpose
    rowWrapper := rowPool.Get().(*Row)
    rowWrapper.SetContext(d.Context())
    rowWrapper.init(row, query, args)

    return rowWrapper
}

// Query perform query using a default timeout context.
func (d *Db) Query(query string, args ...interface{}) *sql.Rows {
    ctx, _ := context.WithTimeout(context.Background(), DefaultDbTimeout)
    return d.QueryContext(ctx, query, args...)
}

// QueryContext perform query using a specified context.
func (d *Db) QueryContext(ctx context.Context, query string, args ...interface{}) *sql.Rows {
    start := time.Now()
    defer func() {
        elapse := time.Since(start)
        d.profileAdd(nil, "Db.Query", elapse, query, args...)

        if elapse >= d.client.SlowLogTime() && d.client.SlowLogTime() > 0 {
            d.Context().Warn("Db.Query slow, elapse:%s, query:%s, args:%v", elapse, query, args)
        }
    }()

    rows, err := d.GetDb(d.master).QueryContext(ctx, query, args...)

    if err != nil {
        d.Context().Error("Db.Query error, %s, query:%s, args:%v", err.Error(), query, args)
        return nil
    }

    return rows
}

// Exec perform exec using a default timeout context.
func (d *Db) Exec(query string, args ...interface{}) sql.Result {
    ctx, _ := context.WithTimeout(context.Background(), DefaultDbTimeout)
    return d.ExecContext(ctx, query, args...)
}

// ExecContext perform exec using a specified context.
func (d *Db) ExecContext(ctx context.Context, query string, args ...interface{}) sql.Result {
    start := time.Now()
    defer func() {
        elapse := time.Since(start)
        d.profileAdd(nil, "Db.Exec", elapse, query, args...)

        if elapse >= d.client.SlowLogTime() && d.client.SlowLogTime() > 0 {
            d.Context().Warn("Db.Exec slow, elapse:%s, query:%s, args:%v", elapse, query, args)
        }
    }()

    res, err := d.GetDb(true).ExecContext(ctx, query, args...)

    if err != nil {
        d.Context().Error("Db.Exec error, %s, query:%s, args:%v", err.Error(), query, args)
        return nil
    }

    return res
}

// Prepare creates a prepared statement for later queries or executions,
// the Close method must be called by caller.
func (d *Db) Prepare(query string) IStmt {
    ctx, _ := context.WithTimeout(context.Background(), DefaultDbTimeout)
    return d.PrepareContext(ctx, query)
}

// PrepareContext creates a prepared statement for later queries or executions,
// the Close method must be called by caller.
func (d *Db) PrepareContext(ctx context.Context, query string) IStmt {

    master, pos := true, strings.IndexByte(query, ' ')
    if pos != -1 && strings.ToUpper(query[:pos]) == "SELECT" {
        master = false
    }

    stmt, err := d.GetDb(master).PrepareContext(ctx, query)

    if err != nil {
        d.Context().Error("Db.Prepare error, %s, query:%s", err.Error(), query)
        return nil
    }

    // wrap stmt for profile purpose
    stmtWrapper := stmtPool.Get().(*Stmt)
    stmtWrapper.SetContext(d.Context())
    stmtWrapper.init(stmt, d.client, query)

    return stmtWrapper
}

type base struct {
    pgo2.Object
    client *db.Client
}

func (b *base) profileAdd(context iface.IContext, key string, elapse time.Duration, query string, args ...interface{}) {
    newKey := key
    if b.client.SqlLog() == true {
        if len(args) > 0 {
            for _, v := range args {
                query = strings.Replace(query, "?", util.ToString(v), 1)
            }
        }

        newKey = key + "(" + query + ")"
    }
    if context != nil {
        context.ProfileAdd(newKey, elapse)
    } else {
        b.Context().ProfileAdd(newKey, elapse)
    }

}

// Tx wrapper for sql.Tx
type Tx struct {
    tx *sql.Tx
    base
}

func (t *Tx) init(tx *sql.Tx, client *db.Client) {
    t.tx = tx
    t.client = client
}

// QueryOne perform one row query using a default timeout context,
// and always returns a non-nil value, Errors are deferred until
// Row's Scan method is called.
func (t *Tx) QueryOne(query string, args ...interface{}) IRow {
    ctx, _ := context.WithTimeout(context.Background(), DefaultDbTimeout)
    return t.QueryOneContext(ctx, query, args...)
}

// QueryOneContext perform one row query using a specified context,
// and always returns a non-nil value, Errors are deferred until
// Row's Scan method is called.
func (t *Tx) QueryOneContext(ctx context.Context, query string, args ...interface{}) IRow {
    start := time.Now()
    defer func() {
        elapse := time.Since(start)
        t.profileAdd(nil, "Tx.QueryOne", elapse, query, args...)

        if elapse >= t.client.SlowLogTime() && t.client.SlowLogTime() > 0 {
            t.Context().Warn("Tx.QueryOne slow, elapse:%s, query:%s, args:%v", elapse, query, args)
        }
    }()
    row := t.tx.QueryRowContext(ctx, query, args...)
    // wrap row for profile purpose
    rowWrapper := rowPool.Get().(*Row)
    rowWrapper.SetContext(t.Context())
    rowWrapper.init(row, query, args)

    return rowWrapper
}

// Query perform query using a default timeout context.
func (t *Tx) Query(query string, args ...interface{}) *sql.Rows {
    ctx, _ := context.WithTimeout(context.Background(), DefaultDbTimeout)
    return t.QueryContext(ctx, query, args...)
}

// QueryContext perform query using a specified context.
func (t *Tx) QueryContext(ctx context.Context, query string, args ...interface{}) *sql.Rows {
    start := time.Now()
    defer func() {
        elapse := time.Since(start)
        t.profileAdd(nil, "Db.Query", elapse, query, args...)

        if elapse >= t.client.SlowLogTime() && t.client.SlowLogTime() > 0 {
            t.Context().Warn("Db.Query slow, elapse:%s, query:%s, args:%v", elapse, query, args)
        }
    }()

    rows, err := t.tx.QueryContext(ctx, query, args...)

    if err != nil {
        t.Context().Error("Db.Query error, %s, query:%s, args:%v", err.Error(), query, args)
        return nil
    }

    return rows
}

// Exec perform exec using a default timeout context.
func (t *Tx) Exec(query string, args ...interface{}) sql.Result {
    ctx, _ := context.WithTimeout(context.Background(), DefaultDbTimeout)
    return t.ExecContext(ctx, query, args...)
}

// ExecContext perform exec using a specified context.
func (t *Tx) ExecContext(ctx context.Context, query string, args ...interface{}) sql.Result {
    start := time.Now()
    defer func() {
        elapse := time.Since(start)
        t.profileAdd(nil, "Db.Exec", elapse, query, args...)

        if elapse >= t.client.SlowLogTime() && t.client.SlowLogTime() > 0 {
            t.Context().Warn("Db.Exec slow, elapse:%s, query:%s, args:%v", elapse, query, args)
        }
    }()

    res, err := t.tx.ExecContext(ctx, query, args...)

    if err != nil {
        t.Context().Error("Db.Exec error, %s, query:%s, args:%v", err.Error(), query, args)
        return nil
    }

    return res
}

// Prepare creates a prepared statement for later queries or executions,
// the Close method must be called by caller.
func (t *Tx) Prepare(query string) IStmt {
    ctx, _ := context.WithTimeout(context.Background(), DefaultDbTimeout)
    return t.PrepareContext(ctx, query)
}

// PrepareContext creates a prepared statement for later queries or executions,
// the Close method must be called by caller.
func (t *Tx) PrepareContext(ctx context.Context, query string) IStmt {

    stmt, err := t.tx.PrepareContext(ctx, query)

    if err != nil {
        t.Context().Error("Db.Prepare error, %s, query:%s", err.Error(), query)
        return nil
    }

    // wrap stmt for profile purpose
    stmtWrapper := stmtPool.Get().(*Stmt)
    stmtWrapper.SetContext(t.Context())
    stmtWrapper.stmt = stmt
    stmtWrapper.client = t.client
    stmtWrapper.query = query

    return stmtWrapper
}

// Commit commit transaction that previously started.
func (t *Tx) Commit() bool {
    if t.tx == nil {
        t.Context().Error("Db.Commit not in transaction")
        return false
    } else {
        if e := t.tx.Commit(); e != nil {
            t.Context().Error("Db.Commit error, " + e.Error())
            return false
        }
        t.tx = nil

        return true
    }
}

// Rollback roll back transaction that previously started.
func (t *Tx) Rollback() bool {
    defer func() {
        t.tx = nil
    }()

    if t.tx == nil {
        t.Context().Error("Db.Rollback not in transaction")
        return false
    } else {
        if e := t.tx.Rollback(); e != nil {
            t.Context().Error("Db.Rollback error, " + e.Error())
            return false
        }
        return true
    }
}

// Row wrapper for sql.Row
type Row struct {
    base
    row   *sql.Row
    query string
    args  []interface{}
}

func (r *Row) init(row *sql.Row, query string, args []interface{}) {
    r.row = row
    r.query = query
    r.args = args
}

func (r *Row) close() {
    r.SetContext(nil)
    r.row = nil
    r.query = ""
    r.args = nil
    rowPool.Put(r)
}

// Scan copies the columns in the current row into the values pointed at by dest.
func (r *Row) Scan(dest ...interface{}) error {
    err := r.row.Scan(dest...)
    if err != nil && err != sql.ErrNoRows {
        r.Context().Error("Db.QueryOne error, %s, query:%s, args:%v", err.Error(), r.query, r.args)
    }

    r.close()
    return err
}

// Stmt wrap sql.Stmt, add context support
type Stmt struct {
    base
    stmt  *sql.Stmt
    query string
}

func (s *Stmt) init(stmt *sql.Stmt, client *db.Client, query string) {
    s.stmt = stmt
    s.client = client
    s.query = query
}

// Close close sql.Stmt and return instance to pool
func (s *Stmt) Close() {
    s.SetContext(nil)
    s.stmt.Close()
    s.stmt = nil
    s.query = ""
    stmtPool.Put(s)
}

// QueryOne perform one row query using a default timeout context,
// and always returns a non-nil value, Errors are deferred until
// Row's Scan method is called.
func (s *Stmt) QueryOne(args ...interface{}) IRow {
    ctx, _ := context.WithTimeout(context.Background(), DefaultDbTimeout)
    return s.QueryOneContext(ctx, args...)
}

// parseArgs get Context
func (s *Stmt) parseArgs(args []interface{}) ([]interface{}, iface.IContext) {
    if len(args) == 0 {
        return args, nil
    }
    var context iface.IContext

    var ctrK int
    for k, iv := range args {
        if v, ok := iv.(iface.IContext); ok {
            context = v
            ctrK = k
            break
        }
    }

    if context == nil {
        return args, nil
    }

    if ctrK == 0 {
        return args[1:], context
    }
    l := len(args) - 1
    ret := make([]interface{}, 0, l)

    ret = append(ret, args[:ctrK]...)
    lastStart := ctrK + 1
    if lastStart > l {
        return ret, context
    }

    ret = append(ret, args[lastStart:]...)

    return ret, context
}

// QueryOneContext perform one row query using a specified context,
// and always returns a non-nil value, Errors are deferred until
// Row's Scan method is called.
func (s *Stmt) QueryOneContext(ctx context.Context, args ...interface{}) IRow {
    start := time.Now()
    args, context := s.parseArgs(args)
    if context == nil {
        context = s.Context()
    }

    defer func() {
        elapse := time.Since(start)
        s.profileAdd(context, "Db.StmtQueryOne", elapse, s.query, args...)

        if elapse >= s.client.SlowLogTime() && s.client.SlowLogTime() > 0 {
            context.Warn("Db.StmtQueryOne slow, elapse:%s, query:%s, args:%v", elapse, s.query, args)
        }
    }()

    row := s.stmt.QueryRowContext(ctx, args...)

    // wrap row for profile purpose
    rowWrapper := rowPool.Get().(*Row)
    rowWrapper.SetContext(context)
    rowWrapper.init(row, s.query, args)

    return rowWrapper
}

// Query perform query using a default timeout context.
func (s *Stmt) Query(args ...interface{}) *sql.Rows {
    ctx, _ := context.WithTimeout(context.Background(), DefaultDbTimeout)
    return s.QueryContext(ctx, args...)
}

// QueryContext perform query using a specified context.
func (s *Stmt) QueryContext(ctx context.Context, args ...interface{}) *sql.Rows {
    args, context := s.parseArgs(args)
    if context == nil {
        context = s.Context()
    }

    start := time.Now()
    defer func() {
        elapse := time.Since(start)
        s.profileAdd(context, "Db.StmtQuery", elapse, s.query, args...)

        if elapse >= s.client.SlowLogTime() && s.client.SlowLogTime() > 0 {
            context.Warn("Db.StmtQuery slow, elapse:%s, query:%s, args:%v", elapse, s.query, args)
        }
    }()

    rows, err := s.stmt.QueryContext(ctx, args...)
    if err != nil {
        context.Error("Db.StmtQuery error, %s, query:%s, args:%v", err.Error(), s.query, args)
        return nil
    }

    return rows
}

// Exec perform exec using a default timeout context.
func (s *Stmt) Exec(args ...interface{}) sql.Result {
    ctx, _ := context.WithTimeout(context.Background(), DefaultDbTimeout)
    return s.ExecContext(ctx, args...)
}

// ExecContext perform exec using a specified context.
func (s *Stmt) ExecContext(ctx context.Context, args ...interface{}) sql.Result {
    args, context := s.parseArgs(args)
    if context == nil {
        context = s.Context()
    }

    start := time.Now()
    defer func() {
        elapse := time.Since(start)
        s.profileAdd(context, "Db.StmtExec", elapse, s.query, args...)

        if elapse >= s.client.SlowLogTime() && s.client.SlowLogTime() > 0 {
            context.Warn("Db.StmtExec slow, elapse:%s, query:%s, args:%v", elapse, s.query, args)
        }
    }()

    res, err := s.stmt.ExecContext(ctx, args...)
    if err != nil {
        context.Error("Db.StmtExec error, %s, query:%s, args:%v", err.Error(), s.query, args)
        return nil
    }

    return res
}
