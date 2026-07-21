package response

import "github.com/gin-gonic/gin"

type Body struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func Success(c *gin.Context, status int, data any) {
	c.JSON(status, Body{
		Code: 0,
		Msg:  "success",
		Data: data,
	})
}

func Error(c *gin.Context, status int, code int, msg string) {
	c.JSON(status, Body{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}
