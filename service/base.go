package service

import (
	"bytes"
	"encoding/json"
	"github.com/emicklei/go-restful"
	"github.com/labstack/gommon/log"
)

type Handler struct {
}

type RespStruct struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Count   int64       `json:"count"`
	Msg     string      `json:"msg"`
	Uuid    string      `json:"uuid"`
	Error   string      `json:"error"`
}

func Finish(response *restful.Response, resp RespStruct) {
	if err := response.WriteEntity(resp); err != nil {
		log.Error("response write Entity err: ", err)
	}
	return
}

func WriteHtml(request *restful.Request, response *restful.Response, resp interface{}) {
	if err := response.WriteEntity(resp); err != nil {
		log.Error("response write Entity err: ", err)
	}
	//log.Infof("uuid:%s ,method:%s ,url:%s , return:%+v",
	//	resp.Uuid, request.Request.Method, request.Request.URL.Path, resp)
	return
}

// InitRequestWithBody
func InitRequestWithBody(request *restful.Request, response *restful.Response, receiveStruct interface{}) (RespStruct, error) {
	buf := new(bytes.Buffer)
	resp := RespStruct{}
	baseInfo := request.Attribute("baseInfo").(*BaseInfo)
	resp.Uuid = baseInfo.Uuid
	_, err := buf.ReadFrom(request.Request.Body)
	if err != nil {
		resp.Error = err.Error()
		Finish(response, resp)
		return resp, err
	}
	reqBytes := buf.Bytes()
	if len(reqBytes) > 0 {
		err = json.Unmarshal(reqBytes, receiveStruct)
		if err != nil {
			resp.Error = err.Error()
			Finish(response, resp)
			return resp, err
		}
	}
	resp.Success = true
	return resp, err
}
