package model

type Cluster struct {
	ID           int    `orm:"column(id);"`
	ClusterName  string `orm:"column(cluster_name);"`
	API          string `orm:"column(api);"`
	Registry     string `orm:"column(registry)"`
	Version      string `orm:"column(version)"`
	NetWorkPlug  string `orm:"column(network_plug)"`
	ETCD         string `orm:"column(etcd)"`
	ControlPlane string `orm:"column(control_plane)"`
	PodDir       string `orm:"column(pod_dir)"`
	ServiceCidr  string `orm:"column(service_cidr)"`
	Status       string `orm:"column(status)"`
}

func (c *Cluster) TableName() string {
	return "cluster"
}

func AddCluster(c *Cluster) (int, error) {
	err := db.Table("cluster").Save(c).Error
	if err != nil {
		return 0, err
	}
	return c.ID, err
}
