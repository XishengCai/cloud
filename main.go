package main

import (
	"cloud/common"
	"cloud/constant"
	"cloud/service"
	"fmt"
	"github.com/urfave/cli"
	"net/http"
	"os"
	"runtime"

	"github.com/emicklei/go-restful"
	"github.com/labstack/gommon/log"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	app := cli.NewApp()
	app.Version = "0.0.1"
	app.Name = "xixi cloud"
	var config string
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config",
			Usage:       "set config path",
			Value:       "",
			Destination: &config,
		},
	}

	app.Action = func(c *cli.Context) error {
		Run(config)
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func Run(config string) {
	cf := common.GetConf(config)
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

	// webservice register to global webservice list
	service.RegisterCluster()
	service.RegisterStatic()

	// webservice add to webservice container
	ar := service.AllRoute{}
	ar.AddAllWebService(wsContainer)

	// start static server
	wsContainer.Handle("/apidocs/",
		http.StripPrefix("/apidocs/",
			http.FileServer(http.Dir(constant.SWAGGER_UI_DIR))))

	fmt.Println("http://localhost:8080/apidocs/?url=http://localhost:8080/download/apidocs.yaml")

	server := &http.Server{Addr: ":" + cf.Port, Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
