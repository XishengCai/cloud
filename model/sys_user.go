package model

import "fmt"

type SysUser struct {
	Id          int       `orm:"column(id)"`
	Username    string    `orm:"column(username);size(32)"`
	Password    string    `orm:"column(password);size(128)"`
	Name        string    `orm:"column(name);size(128);null"`
	MobilePhone string    `orm:"column(mobile_phone);size(32);null"`
	Email       string    `orm:"column(email);size(128);null"`
	Status      string    `orm:"column(status);size(32)"`
	Group       string    `orm:"column(group);size(512);null"`
	UserType    string    `orm:"column(user_type);size(128);null"`
}

func (t *SysUser) TableName() string {
	return "sys_user"
}

func AddSysUser(v SysUser) error {
	if !db.NewRecord(v) {
		return fmt.Errorf("primary key has existed")
	}
	tmp := db.Create(&v)
	return tmp.Error
}

func GetSysUserByName(username string) (sysUser SysUser,  err error) {
	err = db.Where("username = ? and user_type in ('all','interface') " +
		"and status = 0 ", username).First(&sysUser).Error
	return
}
