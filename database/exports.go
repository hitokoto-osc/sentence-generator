package database

import (
	"github.com/hitokoto-osc/hitokoto-sentence-generator/config"
	"github.com/hitokoto-osc/hitokoto-sentence-generator/logging"
	"github.com/upper/db/v4"
)

// Session is upper database session
var Session db.Session

// Connect connect to database
func Connect() (err error) {
	switch config.Database.Driver {
	case "mysql":
		err = mysqlConnect()
	default:
		logging.Logger.Fatal("[database.connect] unsupported database driver.")
	}
	return err
}
