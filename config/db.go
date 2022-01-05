package config

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// configuration for connecting to database
const (
	username string = "fauzil"
	password string = "dbF0r-ummat"
	database string = "quiz-3"
)

// data source
// dsn = "username:password@tcp(hostname)/dbname"
var dsn = fmt.Sprintf("%v:%v@/%v", username, password, database)

func MySQL() (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	return db, nil
}
