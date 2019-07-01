package install

import (
	"log"
)
type Cluster struct{
	Name  string
	MasterHost []string
}

func(c *Cluster) installMaster(){
	for _, mh := range c.MasterHost {
		log.Printf("master host: %s", mh)

	}
}

func(c *Cluster) installSlave(){

}

func(c *Cluster) deleteSlave(){

}

func(c *Cluster) deleteMaster(){

}
