package handler

import (
	"github.com/cloud/common"
	"github.com/cloud/host"
	"github.com/emicklei/go-restful"
)

func (h *Handler)HostAdd(request *restful.Request, response *restful.Response) {
	// 解析
}

func (h *Handler) HostList(request *restful.Request, response *restful.Response) {
	var receiveObject host.Host
	var err error
	resp := h.InitRequestWithBody(request, response, receiveObject)
	receiveObject.Page, receiveObject.PageSize, err = common.GatePage(request)
	if err != nil {
		resp.Error = err
		resp.Success = false
		h.Finish(request, response, resp)
	}
	resp.Data, resp.Count, resp.Error = receiveObject.List()
	h.Finish(request, response, resp)
	return
}

func (h *Handler)HostDelete(request *restful.Request, response *restful.Response) {

}

func (h *Handler) HostUpdate(request *restful.Request, response *restful.Response) {

}
