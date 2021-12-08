package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/mitchellh/mapstructure"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

//InClauseQuery prepare args for In Clause query
//accept query and params
// return query , args and error
func InClauseQuery(query string, args []interface{}) (string, []interface{}, error) {
	q, sqlArgs, err := sqlx.In(query, args) //creates the query string and arguments
	if err != nil {
		return "", nil, err
	}
	//you should check for errors of course
	q = sqlx.Rebind(sqlx.DOLLAR, q)
	return q, sqlArgs, nil
}

//RollbackAndLogError rollback tx with error logging
//accept tx,logger and query
func RollbackAndLogError(tx TxMgr, log *zap.Logger, query string) {
	log.Error(fmt.Sprintf("rolling back tx %s", query))
	if errTx := tx.Rollback(); errTx != nil {
		log.Error(fmt.Sprintf("failed to rollback tx %s", query))
	}
}

//InClauseQueryTwoParams prepare args for In Clause query  with two params (for complex query)
//accept query and params
// return query , args and error
func InClauseQueryTwoParams(query string, argOne, argTwo []interface{}) (string, []interface{}, error) {
	q, sqlArgs, err := sqlx.In(query, argOne, argTwo) //creates the query string and arguments
	if err != nil {
		return "", nil, err
	}
	//you should check for errors of course
	q = sqlx.Rebind(sqlx.DOLLAR, q)
	return q, sqlArgs, nil
}

//ExecuteTx execute query in tx
func (th TxHandler) SeclectQueryTx(ctx context.Context, query string, obj interface{}, params ...interface{}) error {
	tx, err := th.getTxManager(ctx)
	if err != nil {
		return err
	}
	var rows RowMgr
	if len(params) == 0 {
		rows, err = tx.Query(query)
	} else {
		rows, err = tx.Query(query, params...)
	}
	if err != nil {
		err := th.rollBack(ctx)
		if err != nil {
			return err
		}
	}
	return th.processResultSet(ctx, rows, obj)
}

func (th TxHandler) processResultSet(ctx context.Context, rows RowMgr, obj interface{}) error {
	cols, _ := rows.Columns()
	mArr := make([]map[string]interface{}, 0)
	for rows.Next() {
		m := make(map[string]interface{})
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			err := th.rollBack(ctx)
			if err != nil {
				return err
			}
		}
		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}
		mArr = append(mArr, m)
	}
	return mapstructure.Decode(mArr, obj)
}

func (th TxHandler) SeclectQueryInClauseTx(ctx context.Context, obj interface{}, query string, params []interface{}) error {
	if len(params) == 0 {
		err := th.rollBack(ctx)
		if err != nil {
			return err
		}
		return fmt.Errorf("missing query params")
	}
	tx, err := th.getTxManager(ctx)
	if err != nil {
		return err
	}
	q, args, err := InClauseQuery(query, params)
	if err != nil {
		return err
	}
	rows, err := tx.Query(q, args...)
	return th.processResultSet(ctx, rows, obj)
}

//ExecuteTx execute query in tx
func (th TxHandler) ExecuteTx(ctx context.Context, query string, params ...string) (sql.Result, error) {
	tx, err := th.getTxManager(ctx)
	if err != nil {
		return nil, err
	}
	return tx.Exec(query, params)
}

// ExecuteInClauseTx ExecuteTx execute query in tx
func (th TxHandler) ExecuteInClauseTx(ctx context.Context, query string, params []interface{}) (sql.Result, error) {
	tx, err := th.getTxManager(ctx)
	if err != nil {
		return nil, err
	}
	q, args, err := InClauseQuery(query, params)
	if err != nil {
		return nil, err
	}
	rs, err := tx.Exec(q, args...)
	if err != nil {
		err := th.rollBack(ctx)
		if err != nil {
			return nil, err
		}
	}
	return rs, err
}

func NewTxHandler(pgDriver PostgresqlDriver) *TxHandler {
	return &TxHandler{PgDriver: pgDriver}
}

type TxHandler struct {
	PgDriver PostgresqlDriver
}

type Tx struct {
	isRolledBack bool
	txManager    TxMgr
}

func (th *TxHandler) getTxManager(ctx context.Context) (TxMgr, error) {
	contextValue := ctx.Value("tx")
	var ctxn context.Context
	var txData *Tx
	if contextValue == nil {
		ctxn = GetTxContext(ctx)
		contextValue = ctxn.Value("tx")
	}
	txData = contextValue.(*Tx)
	if txData.isRolledBack {
		return nil, fmt.Errorf("tx has been rolled back")
	}
	if txData.txManager == nil {
		tx, err := th.PgDriver.Begin()
		if err != nil {
			return tx, err
		}
		txData.txManager = tx
	}
	return txData.txManager, nil
}
func GetTxContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, "tx", &Tx{})
}

func (th TxHandler) rollBack(ctx context.Context) error {
	tx, err := th.getTxManager(ctx)
	if err != nil {
		return err
	}
	err = tx.Rollback()
	if err != nil {
		return err
	}
	th.markTxAsRollBacked(ctx)
	return nil
}

func (th TxHandler) Commit(ctx context.Context) error {
	tx, err := th.getTxManager(ctx)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (th TxHandler) markTxAsRollBacked(ctx context.Context) {
	contextValue := ctx.Value("tx")
	txData := contextValue.(*Tx)
	if txData.txManager != nil {
		txData.isRolledBack = true
	}
}
