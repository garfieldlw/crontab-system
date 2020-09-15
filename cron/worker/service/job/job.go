package job

import (
	common2 "github.com/garfieldlw/crontab-system/cron/worker/service/common"
	"github.com/garfieldlw/crontab-system/cron/worker/service/handler"
	"github.com/garfieldlw/crontab-system/library/etcd"
	"github.com/garfieldlw/crontab-system/library/log"
	"github.com/garfieldlw/crontab-system/library/utils"
	"github.com/garfieldlw/crontab-system/page/service/common"
	"github.com/garfieldlw/crontab-system/page/service/lock"
	log2 "github.com/garfieldlw/crontab-system/page/service/log"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/copier"
	"go.uber.org/zap"
	"math/rand"
	"sync"
	"time"
)

type DoingJobModel struct {
	FlowId string
	JobId  string
	Cancel context.CancelFunc
}

var (
	doingJobs = &sync.Map{}
)

func DoJobs(ctx context.Context) {
	for {
		select {
		case job := <-common2.JobChan:
			//处理获得到的任务事件
			handleJobEvent(ctx, job)
		}
	}
}

func PushJobChan(job *common.JobInfo) {
	if utils.Contains(common2.JobList, job.JobId) {
		common2.JobChan <- job
	}
}

func handleJobEvent(ctx context.Context, job *common.JobInfo) {
	switch job.JobType {
	case 1:
		//强杀任务
		mapKey := getMapKey(job.FlowId, job.JobId)
		if val, ok := doingJobs.Load(mapKey); ok {
			if val == nil {
				return
			}

			if m, b := val.(*DoingJobModel); b {
				m.Cancel()
			}

			doingJobs.Delete(mapKey)
		}
	default:
		doJob(ctx, job)
	}
}

func getMapKey(flowId, jobId string) string {
	return fmt.Sprintf("%v.%v", flowId, jobId)
}

func doJob(ctx context.Context, info *common.JobInfo) {
	//随机延时
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

	go func(info *common.JobInfo) {
		ctxNew := utils.CopyContext(ctx)

		//抢锁
		var jobLock *lock.Lock
		lockCount := 0
		for {
			if lockCount > 5 {
				log.Warn("get job lock fail", zap.Any("job info", info))
				return
			}

			jobLock = lock.CreateLock("crontab", etcd.PrefixEnumJob, fmt.Sprintf("lock/%v/%v", info.FlowId, info.JobId))
			err := jobLock.Lock()
			if err == nil {
				break
			}

			log.Warn("get job lock fail", zap.Any("job info", info), zap.Error(err))

			if jobLock != nil {
				jobLock.Unlock()
			}

			lockCount = lockCount + 1

			time.Sleep(time.Second * 10)
		}
		defer jobLock.Unlock()

		log.Info("do job, job map")

		ctxNew, ctxCancel := context.WithCancel(ctxNew)

		mapKey := getMapKey(info.FlowId, info.JobId)
		doingJobs.Store(mapKey, &DoingJobModel{
			FlowId: info.FlowId,
			JobId:  info.JobId,
			Cancel: ctxCancel,
		})

		log.Info("do job", zap.Any("job info", info))

		inputData, errInputData := json.Marshal(info.JobInput)
		if errInputData != nil {
			log.Warn("input json marshal error", zap.Any("job info", info), zap.Error(errInputData))
		}

		logEntity, _ := log2.CreateLogJob(ctxNew, 0, 0, info.TraceId, info.FlowId, info.FlowName, info.FlowInfo, info.JobId, info.JobName, info.JobInfo, time.Now().Unix(), 0, string(inputData[:]), "", "", "")

		resultInfo := new(common.JobResultInfo)
		_ = copier.Copy(resultInfo, info)

		defer func() {
			jsonData, _ := json.Marshal(resultInfo)
			_ = etcd.GetDefaultEtcdService().Put("crontab", etcd.PrefixEnumJob, fmt.Sprintf("result/%v/%v", info.FlowId, info.JobId), string(jsonData[:]))
			_ = etcd.GetDefaultEtcdService().Delete("crontab", etcd.PrefixEnumJob, fmt.Sprintf("list/%v/%v", info.FlowId, info.JobId))
			if logEntity != nil && logEntity.Id > 0 {
				outputData, errOutputData := json.Marshal(resultInfo.JobOutput)
				if errOutputData != nil {
					log.Warn("output json marshal error", zap.Any("job info", resultInfo), zap.Error(errOutputData))
				}

				_, _ = log2.UpdateLogJob(ctxNew, logEntity.Id, 0, 0, "", "", "", "", "", "", 0, time.Now().Unix(), "", string(outputData[:]), resultInfo.ErrorMsg, "")
			}
		}()

		//抢到了锁
		result, errResult := handler.Do(ctxNew, info)
		if errResult != nil {
			resultInfo.ErrorMsg = errResult.Error()
		}

		resultInfo.JobOutput = result
	}(info)
}
