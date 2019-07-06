package host

import (
	. "github.com/cloud/common"
	"github.com/cloud/model"
	"github.com/golang/glog"
)

type Host struct {
	IP     string `json:"ip",default:`
	Memory int    `json:"memory"`
	CPU    int    `json:"cpu"`
	Disk   int    `json:"disk"`
	BaseParam
}

func (h *Host) List() ([]model.Host, int64, error) {
	glog.Info("get host list")
	offset := h.Page * h.PageSize
	return model.GetHostList(offset, h.PageSize, "")
}

func (h *Host) Add() {

}

func (h *Host) Delete() {

}

func (h *Host) Update() {

}
