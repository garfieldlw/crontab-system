package flow

import (
	"github.com/garfieldlw/crontab-system/app/controller"
	"github.com/garfieldlw/crontab-system/library/format/format-api"
	"github.com/garfieldlw/crontab-system/page/service/common"
	"github.com/garfieldlw/crontab-system/page/service/flow"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	*controller.BaseController
}

func (o *Controller) List(c *gin.Context) {
	var reqInfo common.CrontabFlowListInputModel
	errReq := o.BindRequestJsonBody(c, &reqInfo)
	if errReq != nil {
		res, _ := format_api.ApiJsonWithError(format_api.ParameterError, errReq.Error(), "")
		c.JSON(http.StatusOK, res)
		return
	}

	result, err := flow.List(c.Request.Context(), reqInfo.Id, reqInfo.Name, reqInfo.FlowType, reqInfo.Status, reqInfo.SortValue, reqInfo.Offset, reqInfo.Limit)
	if err != nil {
		res, _ := format_api.ApiJsonWithError(format_api.ServerError, err.Error(), "")
		c.JSON(http.StatusOK, res)
		return
	}

	res, _ := format_api.ApiJsonWithError(format_api.CodeSuccess, "", result)
	c.JSON(http.StatusOK, res)
	return
}

func (o *Controller) Detail(c *gin.Context) {
	var reqInfo common.CrontabFlowDetailInputModel
	errReq := o.BindRequestJsonBody(c, &reqInfo)
	if errReq != nil {
		res, _ := format_api.ApiJsonWithError(format_api.ParameterError, errReq.Error(), "")
		c.JSON(http.StatusOK, res)
		return
	}

	result, err := flow.DetailById(c.Request.Context(), reqInfo.Id)
	if err != nil {
		res, _ := format_api.ApiJsonWithError(format_api.ServerError, err.Error(), "")
		c.JSON(http.StatusOK, res)
		return
	}

	res, _ := format_api.ApiJsonWithError(format_api.CodeSuccess, "", result)
	c.JSON(http.StatusOK, res)
	return
}

func (o *Controller) Do(c *gin.Context) {
	var reqInfo common.FlowDoInputModel
	errReq := o.BindRequestJsonBody(c, &reqInfo)
	if errReq != nil {
		res, _ := format_api.ApiJsonWithError(format_api.ParameterError, errReq.Error(), "")
		c.JSON(http.StatusOK, res)
		return
	}

	err := flow.Do(c.Request.Context(), reqInfo.FlowId, reqInfo.DoForce, reqInfo.Date)
	if err != nil {
		res, _ := format_api.ApiJsonWithError(format_api.ServerError, err.Error(), "")
		c.JSON(http.StatusOK, res)
		return
	}

	res, _ := format_api.ApiJsonWithError(format_api.CodeSuccess, "", "")
	c.JSON(http.StatusOK, res)
	return
}
