package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Gin struct {
	C *gin.Context
}

const (
	RequestResList string = "requestResList"
)

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
		Status: httpCode,
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
		Status: 0,
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

func NewResponseDateItemWithStatus(err error, status int) Response {
	return Response{
		ResMsg: err.Error(),
		Status: status,
		Data: Data{
			Total: 0,
			Items: []interface{}{},
		},
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
func HandleDataAndError(ctx *gin.Context, eventCode int, data interface{}, err error) {
	if err != nil {
		ctx.JSON(http.StatusOK, NewResponseWithStatus(err, eventCode))
		return
	}
	ctx.JSON(http.StatusOK, NewResponse(data))
}

func HandleDataItemAndError(ctx *gin.Context, eventCode int, data interface{}, err error) {
	if err != nil {
		ctx.JSON(http.StatusOK, NewResponseDateItemWithStatus(err, eventCode))
		return
	}
	ctx.JSON(http.StatusOK, NewResponse(data))
}

// HandleBatchOperationResult handle batch operation result
func HandleBatchOperationResult(ctx *gin.Context, errCode int, bor BatchOperationResp) {
	if bor != nil {
		ctx.JSON(http.StatusOK, Response{
			Status: errCode,
			Data:   bor,
		})
		return
	}
	ctx.JSON(http.StatusOK, NewResponse(nil))
}