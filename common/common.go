package common

import (
	"github.com/cloud/constant"
	"github.com/emicklei/go-restful"
	"github.com/labstack/gommon/log"
	"strconv"
	"strings"
)

type BaseParam struct {
	Page      int      `json:"page"`
	PageSize  int      `json:"page_size"`
	Condition []string `json:"condition"`
}



func GatePage(request *restful.Request) (page int, pageSize int, err error) {
	pageT := strings.TrimSpace(request.QueryParameter(constant.PAGE))
	pageSizeT := strings.TrimSpace(request.QueryParameter(constant.PAGE_SIZE))
	log.Infof("page: %s, pageSize: %s", pageT, pageSizeT)
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


