package api

import (
	"github.com/gin-gonic/gin"
	"github.com/xkenmon/maiev/dto"
)

func checkErrAndWrite(c *gin.Context, err error, code int, msg string) bool {
	if err != nil {
		c.JSON(code, dto.ApiMessage{
			Code: code,
			Msg:  msg,
		})
		return true
	}
	return false
}
func writeMsg(c *gin.Context, code int, msg string) {
	c.JSON(code, dto.ApiMessage{
		Code: code,
		Msg:  msg,
	})
}
