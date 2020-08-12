package models

import (
	"database/sql"
	_"github.com/mattn/go-sqlite3"
)
type SQLDB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}
var Db *sql.DB

func Init(dataSource string) error{
	var err error
	Db,err = sql.Open("sqlite3",dataSource)
	if err != nil {
		return err
	}
	return nil
}