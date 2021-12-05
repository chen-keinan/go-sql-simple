package examples

import (
	"context"
	"fmt"
	"github.com/chen-keinan/go-sql-simple/pkg/db"
)

func main() {
	connector := db.NewConnector("user", "password", "5432", "db", "host", "postgres")
	driver, err := db.NewPGDriver(connector)
	if err != nil {
		panic(err)
	}
	pgMgr := db.NewPostgresqlMgr(driver)
	cbg := db.GetTxContext(context.Background())
	th := db.NewTxHandler(pgMgr)
	res, err := th.ExecuteTx(cbg, "update users set email = ? where name = ? ", "test@gmail.com", "david")
	if err != nil {
		panic(err)
	}
	th.Commit(cbg)
	fmt.Println(res)
}
