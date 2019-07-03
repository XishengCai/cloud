package model

type Host struct{
	ID   int     `gorm:"column(id);size(int)"`
	Name string  `gorm:"column(name);size(80)"`
	IP   string  `gorm:"column(ip);size(80)"`
	Password string `gorm:"column(password);size(80)"`
}

func(t *Host) TableName()string{
	return "host"
}

func GetHostList(){

}