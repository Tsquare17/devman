package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func CreateDatabase(name, password string) bool {
	conn := conn(password)

	_, err := conn.Exec("CREATE DATABASE " + name)
	if err != nil {
		return false
	}

	return true
}

func DeleteDatabase(name, password string) bool {
	conn := conn(password)

	_, err := conn.Exec("DROP DATABASE " + name)
	if err != nil {
		return false
	}

	return true
}

func conn(pass string) *sql.DB {
	conn, err := sql.Open("mysql", "root:" + pass + "@tcp(127.0.0.1:3306)/")
	if err != nil {
		panic(err)
	}

	return conn
}
