package utils

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

type Database struct {
	Connection *sql.DB
	User       string
	Passwd     string
	Net        string
	Host       string
	Port       string
	DBName     string
}

var Db Database

func (db *Database) OpenDbConnection(user string, passwd string, net string, host string, port string, dbname string) (string, error) {
	dbConfig := mysql.Config{
		User:   user,
		Passwd: passwd,
		Net:    net,
		Addr:   fmt.Sprintf("%s:%s", host, port),
		DBName: dbname,
	}

	var err error
	db.Connection, err = sql.Open("mysql", dbConfig.FormatDSN())

	if err != nil {
		return "", fmt.Errorf("failed to open the connection with database: %v", err.Error())
	}

	err = db.Connection.Ping()
	if err != nil {
		return "", fmt.Errorf("failed to connect to the database: %v", err.Error())
	}

	db.User = user
	db.Passwd = passwd
	db.Net = net
	db.Host = host
	db.Port = port
	db.DBName = dbname

	return fmt.Sprintf("connected to dbms server: %v:%v", host, port), nil
}
