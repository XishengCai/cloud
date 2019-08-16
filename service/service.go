package service

import (
	"fmt"
	"cloud/model"
	"cloud/util"
	"github.com/emicklei/go-restful"
	"github.com/labstack/gommon/log"
	"net/http"
	"strings"
)

var gWebServiceList []*restful.WebService

func Register(v *restful.WebService) {
	fmt.Println("register route:  ", v.RootPath())
	if v.RootPath() == "/" {
		gWebServiceList = append(gWebServiceList, v)
	}else{
		gWebServiceList = append(gWebServiceList, v.Filter(ContainerFilter))
	}

}

type AllRoute struct {
}

type BaseInfo struct {
	Uuid     string `json:"uuid"`
	Username string `json:"username"`
	UserType string `json:"user_type"`
}

func (b *BaseInfo) CheckAuthentication(request *restful.Request) bool {
	if request.Request.Method == http.MethodOptions {
		return true
	}
	token := request.Request.Header.Get("token")
	if strings.TrimSpace(token) != "" {
		b.UserType = ""
		username := GetTokenFromRedis(token)
		log.Infof("username: %s, token: %s", b.Username, token)
		if username == b.Username {
			return true
		}
	}
	b.UserType = "interface"
	username := strings.TrimSpace(request.Request.Header.Get("username"))
	password := strings.TrimSpace(request.Request.Header.Get("password"))
	if username == "" || password == "" {
		return false
	}
	b.Username = username
	return CheckPassword(username, password)
}

func (a AllRoute) AddAllWebService(container *restful.Container) {
	for _, ws := range gWebServiceList {
		fmt.Println("Add WebService ", ws.RootPath())
		container.Add(ws)
	}

}

func ContainerFilter(request *restful.Request, response *restful.Response,
	chain *restful.FilterChain) {
	var baseInfo = &BaseInfo{}
	uuid := request.HeaderParameter("uuid")
	username := request.HeaderParameter("username")
	if uuid == "" {
		uuid = util.UuidProvide()
	}
	baseInfo.Uuid = uuid
	baseInfo.Username = username

	var fields = make(map[string]interface{})
	fields["PATH"] = request.Request.URL.Path
	fields["METHOD"] = request.Request.Method
	fields["UUID"] = uuid

	log.Infof("uuid:%s method:%s url:%s",
		uuid, request.Request.Method, request.Request.URL.Path)

	if !baseInfo.CheckAuthentication(request) {
		_ = response.WriteHeaderAndEntity(401,
			fmt.Sprintf("username: %s, auth failed", baseInfo.Username))
		return
	}

	request.SetAttribute("baseInfo", baseInfo)
	chain.ProcessFilter(request, response)

}

func CheckPassword(username, password string) bool {
	v, err := model.GetSysUserByName(username)
	if err != nil {
		log.Errorf("check Password err: %v", err)
		return false
	}
	if v.Password != password {
		return false
	}
	return true
}

func GetTokenFromRedis(token string) (username string) {
	if token == "token" {
		username = "token_test"
	}
	return
}
