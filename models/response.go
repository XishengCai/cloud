package models

type ServiceResponse struct {
	Err    error       `json:"-"`
	Status int         `json:"status"`
	ResMsg string      `json:"resMsg"`
	Data   interface{} `json:"data"`
}
