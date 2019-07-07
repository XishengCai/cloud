package handler

import (
	"github.com/emicklei/go-restful"
)

func init(){
	registerCluster()
}

func registerCluster() {
	hd := &Handler{}
	ws := new(restful.WebService)
	ws.
		Path("/cloud").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML) // you can specify this per route as well

	ws.Route(ws.GET("/host").To(hd.HostList))
	Register(ws)
}