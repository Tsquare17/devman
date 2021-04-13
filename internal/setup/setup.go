package setup

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/tsquare17/devman/internal/database/sqlite"
	"github.com/tsquare17/devman/internal/output"
	"github.com/tsquare17/devman/internal/user"
	"github.com/tsquare17/devman/internal/utils"
	"os"
	osUser "os/user"
)

func Init() {
	username := user.GetName()
	u, err := osUser.Lookup(username)
	if err != nil {
		panic(err)
	}

	var configDir = u.HomeDir + "/.devman"
	if !utils.FileExists(configDir) {
		output.Info("config file no exist")
		err := os.Mkdir(configDir, 0755)
		if err != nil {
			panic(err)
		}
	}

	var dbFile = configDir + "/db.sqlite"
	if !utils.FileExists(dbFile) {
		output.Info("db file no exists")
		file, err := os.Create(dbFile)
		if err != nil {
			panic(err)
		}

		_ = file.Close()

		createTables()
	}
}

func createTables() {
	const query = "CREATE TABLE `sites` (id INTEGER PRIMARY KEY, domain VARCHAR(255), path VARCHAR(255))"
	success := sqlite.Statement(query)

	if success {
		output.Info("DevMan DB created.")
	}
}
