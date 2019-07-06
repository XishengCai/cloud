package model

type Host struct {
	ID       int    `gorm:"column(id);"`
	Name     string `gorm:"column(name);size(80)"`
	IP       string `gorm:"column(ip);size(80)"`
	Password string `gorm:"column(password);size(80)"`
}

func (t *Host) TableName() string {
	return "host"
}

func GetHostList(offset int, limit int, filter string) (hosts []*Host, count int64, err error) {
	err = db.Model(&Host{}).Where("ip like ? ", filter).
		Count(&count).Offset(offset).Limit(limit).Find(&hosts).Error
	return
}
