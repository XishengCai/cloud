package handler

import (
	"github.com/cloud/pkg/model"
	"github.com/emicklei/go-restful"
)

func HostAdd(request *restful.Request, response *restful.Response) {
	// 解析
}

func HostList(request *restful.Request, response *restful.Response) {
	model.GetHostList()
}

func HostDelete(request *restful.Request, response *restful.Response) {

}

func HostUpdate(request *restful.Request, response *restful.Response) {

}