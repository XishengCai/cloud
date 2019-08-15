package service

import (
	"fmt"
	"github.com/cloud/common"
	"github.com/cloud/constant"
	"github.com/cloud/host"
	"github.com/cloud/kubernetes"
	"github.com/emicklei/go-restful"
	"net/http"
	"path"
)

func (h *Handler) HostAdd(request *restful.Request, response *restful.Response) {
	// 解析

}

func (h *Handler) HostList(request *restful.Request, response *restful.Response) {
	var receiveObject host.Host
	var err error
	resp, err := InitRequestWithBody(request, response, receiveObject)
	if err != nil {
		return
	}
	receiveObject.Page, receiveObject.PageSize, err = common.GatePage(request)
	if err != nil {
		resp.Error = err.Error()
		resp.Success = false
		Finish(response, resp)
		return
	}
	resp.Data, resp.Count, err = receiveObject.List()
	if err != nil {
		resp.Error = err.Error()
		resp.Success = false
		Finish(response, resp)
		return
	}
	Finish(response, resp)
	return
}

func (h *Handler) HostDelete(request *restful.Request, response *restful.Response) {

}

func (h *Handler) HostUpdate(request *restful.Request, response *restful.Response) {

}

func (h *Handler) InstallCluster(request *restful.Request, response *restful.Response) {
	var k kubernetes.KubeInstall
	resp, err := InitRequestWithBody(request, response, &k)
	if err != nil {
		return
	}

	if err := k.Install(); err != nil {
		resp.Error = err.Error()
		resp.Success = false
	}
	Finish(response, resp)
	return
}


func Index(request *restful.Request, response *restful.Response) {
	fmt.Println("index")
	WriteHtml(request, response, "This is index")
}

func StaticFromPathParam(req *restful.Request, resp *restful.Response) {
	actual := path.Join(constant.DOWN_LOAD, req.PathParameter("subpath"))
	fmt.Printf("serving %s ... (from %s)\n", actual, req.PathParameter("subpath"))
	http.ServeFile(
		resp.ResponseWriter,
		req.Request,
		actual)
}

func StaticFromQueryParam(req *restful.Request, resp *restful.Response) {
	http.ServeFile(
		resp.ResponseWriter,
		req.Request,
		path.Join(constant.DOWN_LOAD, req.QueryParameter("resource")))
}