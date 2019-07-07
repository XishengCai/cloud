package main

import (
	"log"
	"net/http"
	"runtime"

	_ "github.com/cloud/common"
	"github.com/cloud/handler"
	_ "github.com/cloud/model"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful-swagger12"
)

type Cloud struct {
}

func (cloud Cloud) Register(container *restful.Container) {
	hd := &handler.Handler{}
	ws := new(restful.WebService)
	ws.
		Path("/cloud").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML) // you can specify this per route as well

	ws.Route(ws.GET("/host").To(hd.HostList))
	container.Add(ws)
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	wsContainer := restful.NewContainer()
	cloud := Cloud{}
	cloud.Register(wsContainer)

	// You can install the Swagger Service which provides a nice Web UI on your REST API
	// You need to download the Swagger HTML5 assets and change the FilePath location in the config below.
	// Open http://localhost:8080/apidocs and enter http://localhost:8080/apidocs.json in the api input field.
	config := swagger.Config{
		WebServices:    restful.DefaultContainer.RegisteredWebServices(), // you control what services are visible
		WebServicesUrl: "http://localhost:8081",
		ApiPath:        "/apidocs.json",

		// Optionally, specify where the UI is located
		SwaggerPath:     "/apidocs/",
		SwaggerFilePath: "/Users/emicklei/xProjects/swagger-ui/dist"}
	swagger.RegisterSwaggerService(config, restful.DefaultContainer)

	log.Print("start listening on localhost:8082")
	server := &http.Server{Addr: ":8082", Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
