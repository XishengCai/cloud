package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Gin struct {
	C *gin.Context
}

type ServiceResponser interface {
	checkError() error
	getResp() interface{}
	httpStatus() int
}

type ServiceResponse struct {
	Error  error
	Data   interface{}
	Code   int
	Status int
}

func (r ServiceResponse) checkError() error {
	return r.Error
}

func (r ServiceResponse) getResp() interface{} {
	return r.Data
}

func (r ServiceResponse) httpStatus() int {
	return r.Status
}

type Response struct {
	Status int         `json:"status"`
	ResMsg interface{} `json:"resMsg"`
	Data   interface{} `json:"data"`
}

// BatchOperationResp key requestResList
type BatchOperationResp map[string][]ResponseResult

type ResponseResult struct {
	Id     interface{} `json:"id"`
	ResMsg string      `json:"resMsg"` //执行结果描述
	Status int         `json:"status"` //执行错误码，0为成功
}

type Data struct {
	Total int         `json:"total"`
	Items interface{} `json:"items"`
}

type GitData struct {
	Total int           `json:"total"`
	Items []interface{} `json:"items"`
}

// Response setting gin.JSON
func (g *Gin) Response(httpCode int, msg string, data interface{}) {
	g.C.JSON(httpCode, Response{
		ResMsg: msg,
		Data:   data,
	})
}

// ResponseMsg Response setting gin.JSON
func (g *Gin) ResponseMsg(httpCode, status int, resMsg interface{}, data interface{}) {
	g.C.JSON(httpCode, Response{
		Status: status,
		ResMsg: resMsg,
		Data:   data,
	})
}

func NewResponse(data interface{}) Response {
	if data == nil {
		return Response{
			Status: 0,
			ResMsg: "Success",
			Data: Data{
				Total: 0,
				Items: []interface{}{},
			},
		}
	}

	return Response{
		ResMsg: "Success",
		Data:   data,
	}
}

func NewResponseWithStatus(err error, status int) Response {
	return Response{
		ResMsg: err.Error(),
		Status: status,
	}
}

func HandleError(ctx *gin.Context, errCode int, err error) {
	if err != nil {
		ctx.JSON(http.StatusOK, NewResponseWithStatus(err, errCode))
		return
	}
	ctx.JSON(http.StatusOK, Response{
		Status: 0,
		ResMsg: "Success",
		Data:   nil,
	})
}

// HandleDataAndError 失败事件code
func HandleDataAndError(ctx *gin.Context, resp ServiceResponser) {
	if resp.checkError() != nil {
		ctx.JSON(http.StatusOK, NewResponse(resp.checkError()))
		return
	}
	ctx.JSON(resp.httpStatus(), resp)
}
