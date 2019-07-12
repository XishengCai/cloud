package service

import (
	"github.com/cloud/common"
	"github.com/cloud/host"
	"github.com/cloud/kubernetes"
	"github.com/emicklei/go-restful"
)

func (h *Handler)HostAdd(request *restful.Request, response *restful.Response) {
	// 解析

}

func (h *Handler) HostList(request *restful.Request, response *restful.Response) {
	var receiveObject host.Host
	var err error
	resp, err:= h.InitRequestWithBody(request, response, receiveObject)
	if err != nil{
		return
	}
	receiveObject.Page, receiveObject.PageSize, err = common.GatePage(request)
	if err != nil {
		resp.Error = err.Error()
		resp.Success = false
		h.Finish(request, response, resp)
		return
	}
	resp.Data, resp.Count, err= receiveObject.List()
	if err != nil {
		resp.Error = err.Error()
		resp.Success = false
		h.Finish(request, response, resp)
		return
	}
	h.Finish(request, response, resp)
	return
}

func (h *Handler)HostDelete(request *restful.Request, response *restful.Response) {

}

func (h *Handler) HostUpdate(request *restful.Request, response *restful.Response) {

}


func (h *Handler) InstallCluster(request *restful.Request, response *restful.Response) {
	var k kubernetes.KubeInstall
	resp, err := h.InitRequestWithBody(request, response, &k)
	if err != nil {
		return
	}

	if err:= k.Install(); err!= nil {
		resp.Error = err.Error()
		resp.Success = false
	}
	h.Finish(request, response, resp)
	return
}