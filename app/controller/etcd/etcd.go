package etcd

import (
	"github.com/garfieldlw/crontab-system/app/controller"
	format_api "github.com/garfieldlw/crontab-system/library/format/format-api"
	"github.com/garfieldlw/crontab-system/page/service/common"
	"github.com/garfieldlw/crontab-system/page/service/etcd"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	*controller.BaseController
}

func (o *Controller) Keys(c *gin.Context) {
	var reqInfo common.ETCDKeysInputModel
	errReq := o.BindRequestJsonBody(c, &reqInfo)
	if errReq != nil {
		res, _ := format_api.ApiJsonWithError(format_api.ParameterError, errReq.Error(), "")
		c.JSON(http.StatusOK, res)
		return
	}

	result, err := etcd.GetKeys(c.Request.Context(), reqInfo.Client, reqInfo.Prefix, reqInfo.Key)
	if err != nil {
		res, _ := format_api.ApiJsonWithError(format_api.ServerError, err.Error(), "")
		c.JSON(http.StatusOK, res)
		return
	}

	res, _ := format_api.ApiJsonWithError(format_api.CodeSuccess, "", result)
	c.JSON(http.StatusOK, res)
	return
}
