package models

type Host struct {
	IP       string `form:"ip" binding:"required"`
	Password string `form:"password,default=Root&123" binding:"required"`
	Port     int    `form:"port,default=22" binding:"required"`
	User     string `form:"user,default=root" binding:"required"`
}