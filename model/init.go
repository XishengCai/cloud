package model

import (
	"fmt"
	. "github.com/cloud/constant"
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)
var db *gorm.DB

func init(){
	db = connectMysql()

	autoMigrate()
}

func connectMysql() *gorm.DB{
	glog.Infof("Connect to Mysql:  %s:%s ", MYSQL_HOST, MYSQL_PORT)
	mysql := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		MYSQL_USER, MYSQL_PASSWORD, MYSQL_HOST, MYSQL_PORT, MYSQL_DATABASE)
	db, err := gorm.Open("mysql", mysql)
	if err != nil {
		glog.Errorf("mysql connect err: %v", err)
		panic("连接数据库失败")
	}

	db.LogMode(true)

	return db
}


func autoMigrate(){
	db.AutoMigrate(&Host{})
}
