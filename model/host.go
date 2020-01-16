package model

type Host struct {
	ID         int    `orm:"column(id);"`
	Name       string `orm:"column(name);size(80)"`
	IP         string `orm:"column(ip);size(80)"`
	Password   string `orm:"column(password);size(80)"`
	Port       int    `orm:"column(port);size(80)"`
	User       string `orm:"column(user);size(80)"`
	Memory     int    `orm:"column(memory);size(80)"`
	CPU        int    `orm:"column(cpu);size(80)"`
	Disk       int    `orm:"column(disk);size(80)"`
	InternalIP string `orm:"column(internal_ip);size(80)"`
}

func (t *Host) TableName() string {
	return "host"
}

func GetHostList(offset int, limit int, filter string) (hosts []Host, count int64, err error) {
	filter = "%" + filter + "%"
	err = db.Model(&Host{}).Where("ip like ? ", filter).
		Count(&count).Offset(offset).Limit(limit).Find(&hosts).Error
	return
}

func GetHost(ip string) (host Host, err error) {
	err = db.Model(&Host{}).Where("ip = ? ", ip).First(&host).Error
	return
}

func GetHosts(ip []string) (host []Host, err error) {
	err = db.Model(&Host{}).Where("ip in (?)", ip).First(&host).Error
	return
}
