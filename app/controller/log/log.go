package log

import (
	"github.com/garfieldlw/crontab-system/app/controller"
	format_api "github.com/garfieldlw/crontab-system/library/format/format-api"
	"github.com/garfieldlw/crontab-system/page/service/common"
	"github.com/garfieldlw/crontab-system/page/service/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	*controller.BaseController
}

func (o *Controller) ListFlow(c *gin.Context) {
	var reqInfo common.CrontabLogFlowListInputModel
	errReq := o.BindRequestJsonBody(c, &reqInfo)
	if errReq != nil {
		res, _ := format_api.ApiJsonWithError(format_api.ParameterError, errReq.Error(), "")
		c.JSON(http.StatusOK, res)
		return
	}

	var id int64
	if reqInfo.Id == nil {
		id = 0
	} else {
		id = reqInfo.Id.ToInt64()
	}

	var fatherId int64
	if reqInfo.FatherId == nil {
		fatherId = 0
	} else {
		fatherId = reqInfo.FatherId.ToInt64()
	}

	result, err := log.ListLogFlow(c.Request.Context(), id, fatherId, reqInfo.WorkIp, reqInfo.FlowId, reqInfo.FlowName, reqInfo.StartTime, reqInfo.EndTime, reqInfo.SortValue, reqInfo.Offset, reqInfo.Limit)
	if err != nil {
		res, _ := format_api.ApiJsonWithError(format_api.ServerError, err.Error(), "")
		c.JSON(http.StatusOK, res)
		return
	}

	res, _ := format_api.ApiJsonWithError(format_api.CodeSuccess, "", result)
	c.JSON(http.StatusOK, res)
	return
}

func (o *Controller) ListJob(c *gin.Context) {
	var reqInfo common.CrontabLogJobListInputModel
	errReq := o.BindRequestJsonBody(c, &reqInfo)
	if errReq != nil {
		res, _ := format_api.ApiJsonWithError(format_api.ParameterError, errReq.Error(), "")
		c.JSON(http.StatusOK, res)
		return
	}

	var id int64
	if reqInfo.Id == nil {
		id = 0
	} else {
		id = reqInfo.Id.ToInt64()
	}

	var fatherId int64
	if reqInfo.FatherId == nil {
		fatherId = 0
	} else {
		fatherId = reqInfo.FatherId.ToInt64()
	}

	var traceId int64
	if reqInfo.TraceId == nil {
		traceId = 0
	} else {
		traceId = reqInfo.TraceId.ToInt64()
	}

	result, err := log.ListLogJob(c.Request.Context(), id, fatherId, traceId, reqInfo.WorkIp, reqInfo.FlowId, reqInfo.JobId, reqInfo.FlowName, reqInfo.JobName, reqInfo.StartTime, reqInfo.EndTime, reqInfo.SortValue, reqInfo.Offset, reqInfo.Limit)
	if err != nil {
		res, _ := format_api.ApiJsonWithError(format_api.ServerError, err.Error(), "")
		c.JSON(http.StatusOK, res)
		return
	}

	res, _ := format_api.ApiJsonWithError(format_api.CodeSuccess, "", result)
	c.JSON(http.StatusOK, res)
	return
}
