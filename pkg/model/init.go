package model

import (
	"fmt"
	"github.com/cloud/pkg"
	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
	"os"
)
var db *gorm.DB

func init(){
	db = connectMysql()
}



func connectMysql() *gorm.DB{
	glog.Infof("Connect to Mysql:  %s:%s ", cloud.MYSQL_HOST, cloud.MYSQL_PORT)
	mysql := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		cloud.MYSQL_USER, cloud.MYSQL_PASSWORD, cloud.MYSQL_HOST, cloud.MYSQL_PORT, cloud.MYSQL_DATABASE)
	db, err := gorm.Open("mysql", mysql)
	if err != nil {
		glog.Errorf("mysql connect err: %v", err)
		os.Exit(1000)
	}
	return db
}