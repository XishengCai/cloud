package main

import (
	"fmt"
	"cloud/common"
	"cloud/constant"
	"cloud/service"
	"net/http"
	"runtime"

	"github.com/labstack/gommon/log"

	_ "cloud/common"
	_ "cloud/model"
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
	ar.AddAllWebService(wsContainer) // 注册所有的路由

	wsContainer.Handle("/apidocs/",
		http.StripPrefix("/apidocs/",
			http.FileServer(http.Dir(constant.SWAGGER_UI_DIR)))) //静态文件服务器

	fmt.Println("http://localhost:8080/apidocs/?url=http://localhost:8080/download/apidocs.json")

	config := common.GetConf()
	server := &http.Server{Addr: ":" + config.Port, Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
