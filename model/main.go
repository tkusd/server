package model

import (
	"database/sql"
	"log"
	"strings"

	"path/filepath"

	"bitbucket.org/liamstask/goose/lib/goose"
	"github.com/jinzhu/gorm"
	"github.com/tkusd/server/config"
	"github.com/tkusd/server/util"
)

const (
	databaseDir  = "db"
	defaultLimit = 30
	maxLimit     = 100
)

var db gorm.DB

// QueryOption is the query options.
type QueryOption struct {
	Offset int
	Limit  int
	Order  string
}

func (q *QueryOption) ParseOrder() string {
	var arr []string
	split := util.SplitAndTrim(q.Order, ",")

	for _, s := range split {
		if s[0] == '-' {
			arr = append(arr, string(s[1:]), "desc")
		} else {
			arr = append(arr, s)
		}
	}

	return strings.Join(arr, " ")
}

func init() {
	var dbconf *goose.DBConf
	var err error
	env := config.Env

	if env == "" {
		env = "development"
	}

	dbconf, err = goose.NewDBConf(filepath.Join(config.BaseDir, databaseDir), env, "")

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
