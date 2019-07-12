package model

type HostServer struct {
	ID         int    `orm:"column(id)"`
	HostID     int    `orm:"column(host_id)"`
	ServerName string `orm:"column(process)"`
	Port       int    `orm:"column(port)"`
	ClusterID  int    `orm:"column(cluster_id)"`
}

func (h *HostServer) TableName() string {
	return "host_server"
}

func AddHostServer(hs *HostServer) error {
	return db.Create(hs).Error
}
