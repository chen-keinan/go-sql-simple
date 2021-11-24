package db

import (
	"database/sql"
	"fmt"
	"time"
	//no golint
	_ "github.com/lib/pq"
)

//PostgresqlDriver interface
//pgdriver.go
//go:generate mockgen -destination=./mocks/mock_PostgresqlDriver.go -package=mocks . PostgresqlDriver
type PostgresqlDriver interface {
	Close() error
	Begin() (TxMgr, error)
	Query(query string, args ...interface{}) (RowMgr, error)
}

//TxMgr interface
//pgdriver.go
//go:generate mockgen -destination=./mocks/mock_TxMgr.go -package=mocks . TxMgr
type TxMgr interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Commit() error
	Rollback() error
	Query(query string, args ...interface{}) (RowMgr, error)
}

//RowMgr interface
//pgdriver.go
//go:generate mockgen -destination=./mocks/mock_RowMgr.go -package=mocks . RowMgr
type RowMgr interface {
	Scan(dest ...interface{}) error
	Next() bool
	Close() error
	Columns() ([]string, error)
}

//Rows encapsulate sql rows object
type Rows struct {
	rowMgr  *sql.Rows
}

//SQLTx encapsulate sql tx object
type SQLTx struct {
	Tx *sql.Tx
}

//PostgresqlMgr postgresql mgr
type PostgresqlMgr struct {
	PgPool *sql.DB
}

//NewPostgresqlMgr return new postgresql driver instance
func NewPostgresqlMgr(pgDriver *sql.DB) PostgresqlDriver {
	return &PostgresqlMgr{PgPool: pgDriver}
}

//NewPGDriver create new postgresql driver instance
// accept config and logger
// return instance of sql db
func NewPGDriver(c Connector) (*sql.DB,error) {
	// open database
	db, err := connectToPostgresqlWithRetries(c)
	if err != nil {
		return nil,err
	}
 	return db,nil
}

func connectToPostgresqlWithRetries(c Connector) (*sql.DB, error) {
	var err error
	var db *sql.DB
	for i := 0; i < 25; i++ {
		// open connection to postgresql
		psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
			"password=%s dbname=%s sslmode=disable",
			c.host, c.port, c.user, c.password, c.db)
		if db, err = sql.Open(c.sqlType,psqlInfo); err != nil {
			err = fmt.Errorf("failed to connect to pg db with connnection url %s:%s/%s error:%s", c.host, c.port, c.db, err.Error())
			time.Sleep(time.Second * 4)
			continue
		}
		// check db connection
		if err = db.Ping(); err != nil {
			err = fmt.Errorf("failed to check connection to pg db with connnection url %s:%s/%s : error:%s", c.host, c.port, c.db, err.Error())
			time.Sleep(time.Second * 4)
			continue
		}
		// succeed to open connection and test it
		return db, nil
	}
	// failed to open connection / test it 5 times
	return nil, err
}

//Close close pg pool
// return error
func (pdi PostgresqlMgr) Close() error {
	return pdi.PgPool.Close()
}

//Begin sql tx
// return sql tx and error
func (pdi PostgresqlMgr) Begin() (TxMgr, error) {
	tx, err := pdi.PgPool.Begin()
	return SQLTx{Tx: tx}, err
}

//Query sql query
//accept query and args
//return rows and error
func (pdi PostgresqlMgr) Query(query string, args ...interface{}) (RowMgr, error) {
	return pdi.PgPool.Query(query, args...)
}

//Exec execute pg query (insert / update / delete )
// accept query and params
// return sql exec results and error
func (stx SQLTx) Exec(query string, args ...interface{}) (sql.Result, error) {
	return stx.Tx.Exec(query, args...)
}

//Commit commmit sql tx
// return error
func (stx SQLTx) Commit() error {
	return stx.Tx.Commit()
}

//Rollback rollback sql tx
// return error
func (stx SQLTx) Rollback() error {
	return stx.Tx.Rollback()
}

//Query execute seletc query tx
// acceprt select query and args
// return sql rows and error
func (stx SQLTx) Query(query string, args ...interface{}) (RowMgr, error) {
	rows, err := stx.Tx.Query(query, args...)
	return Rows{rowMgr: rows}, err
}

//Scan scan query results
//accept columns variables as pointers
// return error
func (r Rows) Scan(params ...interface{}) error {
	return r.rowMgr.Scan(params...)
}

//Scan scan query results
//accept columns variables as pointers
// return error
func (r Rows) Columns() ([]string, error) {
	return r.rowMgr.Columns()
}

//Next check if there are additional rows
//return bool
func (r Rows) Next() bool {
	return r.rowMgr.Next()
}

//Close close sql rows
// return error
func (r Rows) Close() error {
	return r.rowMgr.Close()
}
