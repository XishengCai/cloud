package service

import (
	"github.com/emicklei/go-restful"
)


func RegisterCluster() {
	hd := &Handler{}
	ws := new(restful.WebService)
	ws.
		Path("/cluster").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML) // you can specify this per route as well

	ws.Route(ws.GET("/host").To(hd.HostList))

	ws.Route(ws.POST("/install").To(hd.InstallCluster))

	Register(ws)
}
