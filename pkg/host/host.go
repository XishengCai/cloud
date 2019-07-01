package host

type Host struct {
	IP string `json:ip`
	Memory int  `json:memory`
	CPU    int  `json:cpu`
	Disk   int  `json:disk`
}

func (host *Host) list(){

}

func (host *Host) add(){

}

func (host *Host) delete(){

}

func (hsot *Host) update(){

}