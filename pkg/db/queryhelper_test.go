package db

import (
	"context"
	"github.com/chen-keinan/go-sql-simple/pkg/db/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInClauseQuery(t *testing.T) {
	s, param, err := InClauseQuery("select * from test where key = ?", []interface{}{"aaa"})
	assert.NoError(t, err)
	assert.True(t, s == "select * from test where key = $1")
	assert.True(t, param[0].(string) == "aaa")
}

func TestInClauseQueryWithError(t *testing.T) {
	_, _, err := InClauseQuery("select * from test where key = ?", nil)
	assert.Error(t, err)
}

func TestInClauseQueryTwoParams(t *testing.T) {
	s, param, err := InClauseQueryTwoParams("select * from test where key = ? and p=?", []interface{}{"aaa"}, []interface{}{"bbb"})
	assert.NoError(t, err)
	assert.True(t, s == "select * from test where key = $1 and p=$2")
	assert.True(t, param[0].(string) == "aaa")
	assert.True(t, param[1].(string) == "bbb")
}

func TestSelectQueryInTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	postgresqlMgr := mocks.NewMockPostgresqlDriver(ctrl)
	postgresqlMgr.EXPECT().Begin().Return(SQLTx{}, nil)
	th := TxHandler{PgDriver: postgresqlMgr}
	cbg := GetTxContext(context.Background())
	tc, err := th.getTxManager(cbg)
	if err != nil {
		t.Errorf("failed")
	}
	assert.NotNilf(t, tc, "tx manager is nil")
}
