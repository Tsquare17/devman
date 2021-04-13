package sqlite

import (
	"database/sql"
	"github.com/tsquare17/devman/internal/user"
	osUser "os/user"
)

func Statement(statement string) bool {
	dbFile := getDbFile()

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		panic(err)
	}

	query, err := db.Prepare(statement)
	if err != nil {
		panic(err)
	}

	_, err = query.Exec()
	if err != nil {
		panic(err)
	}

	err = db.Close()
	if err != nil {
		panic(err)
	}

	return true
}

func getDbFile() string {
	username := user.GetName()
	u, err := osUser.Lookup(username)
	if err != nil {
		panic(err)
	}

	return u.HomeDir + "/.devman/db.sqlite"
}
