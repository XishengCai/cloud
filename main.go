package main

import (
	"fmt"
	"net/http"

	"cloud/pkg/client/mongodb"
	"cloud/pkg/setting"
	"cloud/routers"
	"cloud/service/kubernetes"

	"github.com/gin-gonic/gin"
	"k8s.io/klog"
)

func init() {
	setting.SetUp("")
	mongodb.InitMongoDB()
	kubernetes.SetUpJob()
}

// @title Swagger Example API
// @version 2.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host xisheng.vip:8081
// @BasePath /api
// @query.collection.format multi

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @securitydefinitions.oauth2.application OAuth2Application
// @tokenUrl https://example.com/oauth/token
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.implicit OAuth2Implicit
// @authorizationurl https://example.com/oauth/authorize
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.password OAuth2Password
// @tokenUrl https://example.com/oauth/token
// @scope.read Grants read access
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.accessCode OAuth2AccessCode
// @tokenUrl https://example.com/oauth/token
// @authorizationurl https://example.com/oauth/authorize
// @scope.admin Grants read and write access to administrative information

// @x-extension-openapi {"example": "value on a json format"}
func main() {

	gin.SetMode(setting.EnvConfig.RunMode)

	routersInit := routers.InitRouter()
	readTimeout := setting.HttpSetting.ReadTimeout
	writeTimeout := setting.HttpSetting.WriteTimeout
	maxHeaderBytes := 1 << 20

	endpoint := fmt.Sprintf(":%s", setting.HttpSetting.HttpPort)
	server := &http.Server{
		Addr:           endpoint,
		Handler:        routersInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	if err := server.ListenAndServe(); err != nil {
		klog.Error(err)
	}
}
