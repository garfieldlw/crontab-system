package controller

import (
	"github.com/gin-gonic/gin"
)

type BaseController struct {
}

func (base *BaseController) BindRequestJsonBody(c *gin.Context, obj interface{}) error {
	errReq := c.ShouldBind(obj)
	if errReq != nil {
		return errReq
	}
	return nil
}

func (base *BaseController) BindRequestUri(c *gin.Context, obj interface{}) error {
	errReq := c.ShouldBindUri(obj)
	if errReq != nil {
		return errReq
	}
	return nil
}
