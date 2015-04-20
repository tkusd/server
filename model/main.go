package model

import (
	"log"

	"path/filepath"

	"os"

	"bitbucket.org/liamstask/goose/lib/goose"
	"github.com/jinzhu/gorm"
	"github.com/tommy351/app-studio-server/util"
)

const DB_DIR = "db"

var db gorm.DB

func init() {
	var dbconf *goose.DBConf
	var err error
	env := os.Getenv("GO_ENV")

	if env == "" {
		env = "development"
	}

	dbconf, err = goose.NewDBConf(filepath.Join(util.GetBaseDir(), DB_DIR), env, "")

	if err != nil {
		log.Fatal(err)
	}

	db, err = gorm.Open(dbconf.Driver.Name, dbconf.Driver.OpenStr)

	if err != nil {
		log.Fatal(err)
	}
}
