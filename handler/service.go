package handler

import (
	"github.com/cloud/common"
	"github.com/cloud/util"
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"net/http"
)

var gWebServiceList []*restful.WebService = nil

func init() {
	gWebServiceList = make([]*restful.WebService, 0)
}

type AllRoute struct {
}

func (a AllRoute) AddAllWebService(container *restful.Container) {
	for _, ws := range gWebServiceList {
		glog.Infof("AddWebService %s", ws.RootPath())
		container.Add(ws)
	}

}

func Register(v *restful.WebService) {
	gWebServiceList = append(gWebServiceList, v.Filter(ContainerFilter))
}

func ContainerFilter(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	var baseInfo = &common.BaseInfo{}
	uuid := request.HeaderParameter("uuid")
	if uuid == "" {
		uuid = util.UuidProvide()
	}
	baseInfo.Uuid = uuid

	var fields = make(map[string]interface{})
	fields["PATH"] = request.Request.URL.Path
	fields["METHOD"] = request.Request.Method
	fields["UUID"] = baseInfo.Uuid

	glog.Infof("uuid:%s method:%s url:%s",
		uuid, request.Request.Method, request.Request.URL.Path)

	if !baseInfo.CheckAuthentication(request) {
		response.WriteHeaderAndEntity(
			401, struct{Uuid string; UserName string}{baseInfo.Uuid,baseInfo.Username})
	}

	request.SetAttribute("baseInfo", baseInfo)
	chain.ProcessFilter(request, response)

}

