package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/xkenmon/maiev/config"
	"log"
)

var db *gorm.DB

func init() {
	conf := config.GetConfig()
	var err error
	db, err = gorm.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?parseTime=true",
		conf.Get("db.user"),
		conf.Get("db.password"),
		conf.Get("db.host"),
		conf.Get("db.name")))
	if err != nil {
		log.Fatal("can not open database: " + err.Error())
	}
	logMode := conf.GetBool("db.log")
	db.LogMode(logMode)
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		if defaultTableName == "" {
			return defaultTableName
		}
		// gorm表名默认加后缀s,所以去掉最后一个字符s
		return "maiev_" + defaultTableName[0:len(defaultTableName)-1]
	}
}

func GetDB() *gorm.DB {
	return db
}
