package model

type Cluster struct {
	ID          int    `orm:"column(id);"`
	Name        string `orm:"column(name);"`
	API         string `orm:"column(api);"`
	Registry    string `orm:"column(registry)"`
	Version     string `orm:"column(version)"`
	NetWorkPlug string `orm:"column(network_plug)"`
	ETCD        string `orm:"column(etcd)"`
}

func (c *Cluster) TableName() string {
	return "cluster"
}

func AddCluster(c *Cluster) (int, error){
	err := db.Save(c).Error
	return c.ID, err
}
