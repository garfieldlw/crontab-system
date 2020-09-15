package job

import (
	"github.com/garfieldlw/crontab-system/app/controller"
	format_api "github.com/garfieldlw/crontab-system/library/format/format-api"
	"github.com/garfieldlw/crontab-system/page/service/common"
	"github.com/garfieldlw/crontab-system/page/service/job"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	*controller.BaseController
}

func (o *Controller) List(c *gin.Context) {
	var reqInfo common.CrontabJobListInputModel
	errReq := o.BindRequestJsonBody(c, &reqInfo)
	if errReq != nil {
		res, _ := format_api.ApiJsonWithError(format_api.ParameterError, errReq.Error(), "")
		c.JSON(http.StatusOK, res)
		return
	}

	result, err := job.List(c.Request.Context(), reqInfo.Id, reqInfo.JobType, reqInfo.Name, reqInfo.Status, reqInfo.SortValue, reqInfo.Offset, reqInfo.Limit)
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
	var reqInfo common.CrontabJobDetailInputModel
	errReq := o.BindRequestJsonBody(c, &reqInfo)
	if errReq != nil {
		res, _ := format_api.ApiJsonWithError(format_api.ParameterError, errReq.Error(), "")
		c.JSON(http.StatusOK, res)
		return
	}

	result, err := job.DetailById(c.Request.Context(), reqInfo.Id)
	if err != nil {
		res, _ := format_api.ApiJsonWithError(format_api.ServerError, err.Error(), "")
		c.JSON(http.StatusOK, res)
		return
	}

	res, _ := format_api.ApiJsonWithError(format_api.CodeSuccess, "", result)
	c.JSON(http.StatusOK, res)
	return
}

func (o *Controller) Insert(c *gin.Context) {

}

func (o *Controller) Delete(c *gin.Context) {

}

func (o *Controller) Kill(c *gin.Context) {

}

func (o *Controller) Retry(c *gin.Context) {

}
