package model

import (
	"database/sql"
	"log"

	"path/filepath"

	"os"

	"bitbucket.org/liamstask/goose/lib/goose"
	"github.com/jinzhu/gorm"
	"github.com/tkusd/server/util"
)

const (
	databaseDir  = "db"
	defaultLimit = 30
)

var db gorm.DB

// QueryOption is the query options.
type QueryOption struct {
	Offset int
	Limit  int
	Order  string
}

func init() {
	var dbconf *goose.DBConf
	var err error
	env := os.Getenv("GO_ENV")

	if env == "" {
		env = "development"
	}

	dbconf, err = goose.NewDBConf(filepath.Join(util.GetBaseDir(), databaseDir), env, "")

	if err != nil {
		log.Fatal(err)
	}

	db, err = gorm.Open(dbconf.Driver.Name, dbconf.Driver.OpenStr)

	if err != nil {
		log.Fatal(err)
	}
}

func exists(table string, id string) bool {
	var result sql.NullBool
	db.Raw("SELECT exists(SELECT 1 FROM "+table+" WHERE id = ?)", id).Row().Scan(&result)
	return result.Bool
}
