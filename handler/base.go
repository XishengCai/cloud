package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
)

type Handler struct {
}

type RespStruct struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Count      int64       `json:"count"`
	Msg        string      `json:"msg"`
	Uuid       string      `json:"uuid"`
	Error      error       `json:"error"`
}

func (h *Handler) Finish(request *restful.Request, response *restful.Response, resp RespStruct) {
	fmt.Printf("--resp: %+v", resp)
	response.WriteEntity(resp)
	glog.Infof("uuid:%s ,method:%s ,url:%s , return:%+v",
		resp.Uuid, request.Request.Method, request.Request.URL.Path, resp)
	return
}

// InitRequestWithBody
func (h Handler) InitRequestWithBody(request *restful.Request, response *restful.Response, receiveStruct interface{}) (resp RespStruct) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(request.Request.Body)
	if err != nil {
		resp.Error = err
		h.Finish(request, response, resp)
		return
	}
	reqBytes := buf.Bytes()
	if len(reqBytes) > 0 {
		err = json.Unmarshal(reqBytes, receiveStruct)
		if err != nil {
			resp.Error = err
			h.Finish(request, response, resp)
			return
		}
	}
	resp.Success = true
	return
}
