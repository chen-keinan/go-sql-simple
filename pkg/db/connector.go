package db

//Connector hold sql connection information
type Connector struct {
	user     string
	port     string
	password string
	db       string
	host     string
	sqlType  string
}

func NewConnector(user, password, port, db, host, sqlType string) Connector {
	return Connector{user: user, password: password, port: port, db: db, sqlType: sqlType}
}
