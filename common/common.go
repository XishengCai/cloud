package common

import (
	"github.com/cloud/constant"
	"github.com/cloud/model"
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"net/http"
	"strconv"
	"strings"
)

type BaseParam struct {
	Page      int      `json:"page"`
	PageSize  int      `json:"page_size"`
	Condition []string `json:"condition"`
}

type BaseInfo struct{
	Uuid     string  `json:"uuid"`
	Username     string  `json:"user_name"`
	UserType string  `json:"user_type"`
}

func (b *BaseInfo)CheckAuthentication(request *restful.Request) bool{
	if request.Request.Method == http.MethodOptions {
		return true
	}
	token := request.Request.Header.Get("token")
	if strings.TrimSpace(token) != "" {
		b.UserType = ""
		username := GetTokenFromRedis(token)
		if username == b.Username{
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

func GatePage(request *restful.Request) (page int, pageSize int, err error) {
	pageT := strings.TrimSpace(request.QueryParameter(constant.PAGE))
	pageSizeT := strings.TrimSpace(request.QueryParameter(constant.PAGE_SIZE))
	glog.Infof("page: %s, pageSize: %s", pageT, pageSizeT)
	if pageT != "" {
		page, err = strconv.Atoi(pageT)
		if err != nil {
			return
		}
		page = page - 1
	} else {
		page = 0
	}

	if pageSizeT != "" {
		pageSize, err = strconv.Atoi(pageSizeT)
		if err != nil {
			return
		}
	} else {
		pageSize = 8
	}
	return
}

func GetTokenFromRedis(token string) (username string){


}

func CheckPassword(username, password string) bool{
	v, err := model.GetSysUserByName(username)
	if err != nil {
		glog.Errorf("check Password err: %v", err)
		return false
	}
	if v.Password != password{
		return false
	}
	return true
}