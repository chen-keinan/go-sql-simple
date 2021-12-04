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
     connector:=db.NewConnector("user", "password","5432, "db"", "host", "postgres")
     if err != nil{
        panic(err) 	
    }   
    db.NewPGDriver(connector)
    db.NewPostgresqlMgr()
	th := db.NewTxHandler(New)
	cbg := GetTxContext(context.Background())
	txm, err := th.getTxManager(cbg)
    if err != nil{
        panic(err)
    }
    res,err:=th.ExecuteTx(txm,"update users set email = ? where name = ? ","test@gmail.com","david")
    if err != nil{
        return err
    }
    fmt.print(res)
```
