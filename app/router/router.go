package router

import (
	"github.com/garfieldlw/crontab-system/app/controller/etcd"
	"github.com/garfieldlw/crontab-system/app/controller/flow"
	"github.com/garfieldlw/crontab-system/app/controller/job"
	"github.com/garfieldlw/crontab-system/app/controller/log"
	"github.com/garfieldlw/crontab-system/app/controller/worker"
	"github.com/gin-gonic/gin"
)

func LoadRouter(engine *gin.Engine) {

	logController := new(log.Controller)
	groupLog := engine.Group("/api/log")
	groupLog.POST("/flow/list", logController.ListFlow)
	groupLog.POST("/job/list", logController.ListJob)

	workerController := new(worker.Controller)
	groupWorker := engine.Group("/api/worker")
	groupWorker.POST("/list", workerController.List)

	jobController := new(job.Controller)
	groupJob := engine.Group("/api/job")
	groupJob.POST("/list", jobController.List)
	groupJob.POST("/detail", jobController.Detail)
	groupJob.POST("/insert", jobController.Insert)
	groupJob.POST("/delete", jobController.Delete)
	groupJob.POST("/kill", jobController.Kill)
	groupJob.POST("/retry", jobController.Retry)

	flowController := new(flow.Controller)
	groupFlow := engine.Group("/api/flow")
	groupFlow.POST("/list", flowController.List)
	groupFlow.POST("/detail", flowController.Detail)
	groupFlow.POST("/do", flowController.Do)

	etcdController := new(etcd.Controller)
	groupETCD := engine.Group("/api/etcd")
	groupETCD.POST("/list", etcdController.Keys)
}
