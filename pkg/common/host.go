package install

type Host struct {
	ID     int    `json:"id"`
	IP     string `json:"ip"`
	Name   string `json:"name"`
	Remark string `json:"remark"`
}

/*
	检查磁盘剩余空间
*/
func (host *Host)checkDisk() {

}


/*
	添加 host
 */
 func(host *Host) add(){

 }

 /*

  */
func(host *Host) delete(){

}

func (host *Host) getList(){

}