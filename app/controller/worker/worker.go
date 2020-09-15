package worker

import (
	"github.com/garfieldlw/crontab-system/app/controller"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	*controller.BaseController
}

func (o *Controller) List(c *gin.Context) {

}
