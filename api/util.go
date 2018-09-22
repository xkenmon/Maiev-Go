package api

import (
	"github.com/gin-gonic/gin"
	"github.com/xkenmon/maiev/dto"
	"strings"
)

func isSortValid(sort string) bool {
	return strings.EqualFold(sort, "asc") || strings.EqualFold(sort, "desc")
}

func isOrderValid(order, orders string) bool {
	for _, v := range strings.Split(orders, "|") {
		if v == order {
			return true
		}
	}
	return false
}

func writeMsg(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(code, dto.ApiMessage{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}
