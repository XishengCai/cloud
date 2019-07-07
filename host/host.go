package host

import (
	. "github.com/cloud/common"
	"github.com/cloud/model"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/ssh"
)

type Host struct {
	IP        string `json:"ip"`
	Name      string  `json:"name"`
	Password  string `json:"password"`
	User      string  `json:"user"`
	Port      int    `json:"port"`
	Memory    int    `json:"memory"`
	CPU       int    `json:"cpu"`
	Disk      int    `json:"disk"`
	SshClient *ssh.Client
	BaseParam
}

func (h *Host) List() ([]model.Host, int64, error) {
	log.Info("get host list")
	offset := h.Page * h.PageSize
	return model.GetHostList(offset, h.PageSize, "")
}

func (h *Host) Add() {

}

func (h *Host) Delete() {

}

func (h *Host) Update() {

}

func (h *Host) setSshClient() (err error){
	h.SshClient, err = GetSshClient(h.IP, "root", h.Password, h.Port)
	return
}
