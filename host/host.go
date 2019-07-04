package host

import(
	. "github.com/cloud/common"
	"github.com/cloud/model"
	"github.com/golang/glog"
)
type Host struct {
	IP string `json:ip`
	Memory int  `json:memory`
	CPU    int  `json:cpu`
	Disk   int  `json:disk`
	BaseParm
}

func (host *Host) List() ([]*model.Host, int64, error){
	glog.Info("get host list")
	offset := host.Page * host.PageSize
	return model.GetHostList(offset, host.PageSize,"")

}

func (host *Host) Add(){

}

func (host *Host) Delete(){

}

func (hsot *Host) Update(){

}