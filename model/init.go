package model

import (
	"fmt"
	"github.com/cloud/common"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/gommon/log"
)

var db *gorm.DB

func init() {
	db = connectMysql()

	autoMigrate()
}

func connectMysql() *gorm.DB {
	config := common.GetConf()
	mysql := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		config.Server[config.Env].Mysql.User,
		config.Server[config.Env].Mysql.Password,
		config.Server[config.Env].Mysql.IP,
		config.Server[config.Env].Mysql.Port,
		config.Server[config.Env].Mysql.DateBase)
	db, err := gorm.Open("mysql", mysql)
	if err != nil {
		log.Errorf("mysql connect err: %v", err)
		panic("连接数据库失败")
	}

	db.LogMode(true)

	return db
}

func autoMigrate() {
	db.AutoMigrate(&Cluster{})
	db.AutoMigrate(&Host{})
	db.AutoMigrate(&HostServer{})
	db.AutoMigrate(&SysUser{})
}
