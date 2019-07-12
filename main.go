package main

import (
	"fmt"
	"github.com/cloud/service"
	"github.com/labstack/gommon/log"
	"net/http"
	"runtime"

	_ "github.com/cloud/common"
	_ "github.com/cloud/model"
	"github.com/emicklei/go-restful"
)


func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	wsContainer := restful.NewContainer()
	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{"X-My-Header"},
		AllowedHeaders: []string{"Content-Type", "Accept", "x-requested-with", "Token", "token", "X-Host-Override", "Host"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS", "DELETE", "PUT"},
		CookiesAllowed: true,
		Container:      wsContainer}
	wsContainer.Filter(cors.Filter)
	// Add container filter to respond to OPTIONS
	wsContainer.Filter(wsContainer.OPTIONSFilter)

	ar := service.AllRoute{}
	ar.AddAllWebService(wsContainer)

	// You can install the Swagger Service which provides a nice Web UI on your REST API
	// You need to download the Swagger HTML5 assets and change the FilePath location in the config below.
	// Open http://localhost:8080/apidocs and enter http://localhost:8080/apidocs.json in the api input field.
	//config := swagger.Config{
	//	WebServices:    restful.DefaultContainer.RegisteredWebServices(), // you control what services are visible
	//	WebServicesUrl: "http://localhost:8081",
	//	ApiPath:        "/apidocs.json",
	//
	//	// Optionally, specify where the UI is located
	//	SwaggerPath:     "/apidocs/",
	//	SwaggerFilePath: "/Users/emicklei/xProjects/swagger-ui/dist"}
	//swagger.RegisterSwaggerService(config, restful.DefaultContainer)

	fmt.Println("start listening on localhost:8082")
	server := &http.Server{Addr: ":8082", Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
