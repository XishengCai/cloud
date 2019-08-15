package service

import (
	"github.com/emicklei/go-restful"
)

func init(){
	registerStatic()
}

func registerStatic() {
	ws := new(restful.WebService)
	ws.
		Path("/")

	ws.Route(ws.GET("download/{subpath:*}").To(StaticFromPathParam))
	ws.Route(ws.GET("download").To(StaticFromPathParam))
	ws.Route(ws.GET("index").To(Index))

	Register(ws)
}