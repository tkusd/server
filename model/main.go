package model

import (
	"log"

	"path/filepath"

	"bitbucket.org/liamstask/goose/lib/goose"
	"github.com/jinzhu/gorm"
	"github.com/tommy351/app-studio-server/util"
)

var db gorm.DB

func init() {
	var dbconf *goose.DBConf
	var err error

	dbconf, err = goose.NewDBConf(filepath.Join(util.GetBaseDir(), "db"), "development", "")

	if err != nil {
		log.Fatal(err)
	}

	db, err = gorm.Open(dbconf.Driver.Name, dbconf.Driver.OpenStr)

	if err != nil {
		log.Fatal(err)
	}
}
