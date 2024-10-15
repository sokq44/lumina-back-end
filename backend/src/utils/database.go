package utils

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB = nil

func OpenDbConnection(user string, passwd string, net string, host string, port string, dbname string) error {
	dbConfig := mysql.Config{
		User:   user,
		Passwd: passwd,
		Net:    net,
		Addr:   fmt.Sprintf("%s:%s", host, port),
		DBName: dbname,
	}

	var err error
	db, err = sql.Open("mysql", dbConfig.FormatDSN())

	if err != nil {
		return errors.New("failed to open the connection with database")
	}

	err = db.Ping()
	if err != nil {
		return errors.New("Failed to connect to the database: " + err.Error())
	}

	return nil
}

func GetDbConnection() (*sql.DB, error) {
	if db == nil {
		return nil, errors.New("no connection present")
	}

	return db, nil
}
