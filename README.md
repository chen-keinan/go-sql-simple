<br><img src="./pkg/img/go-sql-simple.png" width="300" alt="sql-simple.png logo"><br>

# go-sql-simple

Go-Sql-simple  is an open source sql lib for for making simple sql
operation with transaction whilst avoiding boilte plate code 


* [Installation](#installation)
* [Supported DBs](#supported-dbs)
* [Usage](#usage)

## Installation

```shell
go get github.com/chen-keinan/go-sql-simple
```

## Supported DBs:

- postgresql

## Usage
```go

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
```
