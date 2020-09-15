package schedule

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	common2 "github.com/garfieldlw/crontab-system/cron/master/service/common"
	"github.com/garfieldlw/crontab-system/cron/master/service/cron"
	"github.com/garfieldlw/crontab-system/library/etcd"
	"github.com/garfieldlw/crontab-system/library/log"
	"github.com/garfieldlw/crontab-system/library/message/email"
	"github.com/garfieldlw/crontab-system/page/service/common"
	"github.com/garfieldlw/crontab-system/page/service/flow"
	"github.com/garfieldlw/crontab-system/page/service/job"
	"github.com/garfieldlw/crontab-system/page/service/lock"
	log2 "github.com/garfieldlw/crontab-system/page/service/log"
	"go.uber.org/zap"
	"sort"
	"sync"
	"time"
)

func DoFlow(ctx context.Context) {
	flows, errFlows := getAllFlows(ctx)
	if errFlows != nil {
		log.Panic("get flow fails")
		return
	}

	c := cron.New(cron.WithSeconds(), cron.WithChain(cron.Recover()))

	for _, fl := range flows {
		if fl == nil {
			continue
		}

		switch fl.Id {
		case "Environ":
			{
				_, errFlow := c.AddFunc(fl.Spec, func() { funcEnviron()() })
				if errFlow != nil {
					log.Error("load flow error", zap.Error(errFlow))
				}
			}
		case "EnvironScrapy":
			{
				_, errFlow := c.AddFunc(fl.Spec, func() { funcEnvironScrapy()() })
				if errFlow != nil {
					log.Error("load flow error", zap.Error(errFlow))
				}
			}
		}

	}

	c.Start()

	for {
		select {
		case fl := <-common2.FlowChan:
			//处理获得到的任务事件
			err := handleFlowEvent(ctx, fl)
			if err != nil {
				log.Error("do flow chan error", zap.Any("flow chan info", fl), zap.Error(err))
			}
		}
	}
}

func funcEnviron() func() {
	return func() {
		ctx := context.Background()

		flowInfo, errFlowInfo := flow.DetailById(ctx, "Environ")
		if errFlowInfo != nil {
			return
		}

		err := handleEvent(ctx, &common.FlowInfo{Force: true}, flowInfo)
		if err != nil {
			log.Error("do flow schedule error", zap.Any("flow info", flowInfo), zap.Error(err))
		}
	}
}

func funcEnvironScrapy() func() {
	return func() {
		ctx := context.Background()

		flowInfo, errFlowInfo := flow.DetailById(ctx, "EnvironScrapy")
		if errFlowInfo != nil {
			return
		}

		err := handleEvent(ctx, &common.FlowInfo{Force: true}, flowInfo)
		if err != nil {
			log.Error("do flow schedule error", zap.Any("flow info", flowInfo), zap.Error(err))
		}
	}
}

func getAllFlows(ctx context.Context) ([]*common.CrontabFlowModel, error) {
	entities, errEntities := flow.List(ctx, "", "", 1, 1, "id desc", 0, 1000)
	if errEntities != nil {
		return nil, errEntities
	}

	return entities.Data, nil
}

func PushToFlowChan(flow *common.FlowInfo) {
	common2.FlowChan <- flow
}

func handleFlowEvent(ctx context.Context, fl *common.FlowInfo) error {
	flowInfo, errFlowInfo := flow.DetailById(ctx, fl.FlowId)
	if errFlowInfo != nil {
		return errFlowInfo
	}

	return handleEvent(ctx, fl, flowInfo)
}

func handleEvent(ctx context.Context, fl *common.FlowInfo, flowInfo *common.CrontabFlowModel) error {
	jobLock := lock.CreateLock("crontab", etcd.PrefixEnumFlow, fmt.Sprintf("lock/%v", flowInfo.Id))
	err := jobLock.Lock()
	defer jobLock.Unlock()
	if err != nil {
		return err
	}

	logEntity, errLogEntity := log2.CreateLogFlow(ctx, 0, 0, flowInfo.Id, flowInfo.Name, flowInfo.Info, time.Now().Unix(), 0, "", "", "", "")
	if errLogEntity != nil {
		log.Error("create flow log error", zap.Error(errLogEntity))
	}

	output := ""
	errMsg := ""
	defer func() {
		if len(errMsg) > 0 {
			_ = email.SendMailTLS("username@email.com", fmt.Sprintf("工作流处理失败[%v|%v]", flowInfo.Id, flowInfo.Name), errMsg, "")
		} else {
			_ = email.SendMailTLS("username@email.com", fmt.Sprintf("工作流处理成功[%v|%v]", flowInfo.Id, flowInfo.Name), errMsg, "")
		}

		if logEntity != nil && logEntity.Id > 0 {
			_, errLogEntity := log2.UpdateLogFlow(ctx, logEntity.Id, 0, "", "", "", 0, time.Now().Unix(), "", output, errMsg, "")
			if errLogEntity != nil {
				log.Error("create flow log error", zap.Error(errLogEntity))
			}
		}
	}()

	errDelete := etcd.GetDefaultEtcdService().Delete("crontab", etcd.PrefixEnumFlow, "list/"+flowInfo.Id)
	if errDelete != nil {
		errMsg = errDelete.Error()
		return errDelete
	}

	var traceId int64 = 0
	if logEntity != nil && logEntity.Id > 0 {
		traceId = logEntity.Id
	}

	log.Info(flowInfo.Id, zap.Any("flow info", flowInfo))
	if flowInfo == nil {
		errMsg = "flow info is empty"
		return errors.New(errMsg)
	}

	if fl == nil {
		fl = &common.FlowInfo{}
	}

	if fl.Date == 0 {
		fl.Date = time.Now().Unix()
	}

	var info *common.FlowInfoModel
	errJson := json.Unmarshal([]byte(flowInfo.Info[:]), &info)
	if errJson != nil {
		errMsg = errJson.Error()
		return errJson
	}

	if info == nil || info.Jobs == nil {
		errMsg = "job is empty"
		return errors.New(errMsg)
	}

	var keys []int
	for k := range info.Jobs {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)

	jobResults := make(map[string]*common.JobResultInfo)
	jobOutput := make(map[string]interface{})
	for _, k := range keys {
		jobs := info.Jobs[int32(k)]
		if jobs == nil || len(jobs) == 0 {
			continue
		}

		doForce := false
		if fl.Force {
			doForce = true
		} else {
			if fl.DoForce != nil {
				if _, ok := fl.DoForce[int32(k)]; ok {
					doForce = true
				}
			}
		}

		var wg sync.WaitGroup
		for _, jobId := range jobs {
			if len(jobId) == 0 {
				continue
			}

			//set key
			jobDetail, errJob := job.DetailById(ctx, jobId)
			if errJob != nil {
				errMsg = errJob.Error()
				return errJob
			}

			jobInfo := new(common.JobInfo)
			jobInfo.DoForce = doForce
			jobInfo.TraceId = traceId
			jobInfo.FlowId = flowInfo.Id
			jobInfo.FlowName = flowInfo.Name
			jobInfo.FlowInfo = flowInfo.Info
			jobInfo.FlowTye = flowInfo.FlowType
			jobInfo.JobId = jobDetail.Id
			jobInfo.JobName = jobDetail.Name
			jobInfo.JobType = jobDetail.JobType
			jobInfo.JobInfo = jobDetail.Info
			jobInfo.JobInput = map[string]interface{}{
				"PreJobOutput": jobOutput,
			}
			jobInfo.FlowInput = map[string]interface{}{
				"date": fl.Date,
			}
			jsonData, errJson := json.Marshal(jobInfo)
			if errJson != nil {
				errMsg = errJson.Error()
				return errJson
			}

			errPut := etcd.GetDefaultEtcdService().Put("crontab", etcd.PrefixEnumJob, fmt.Sprintf("list/%v/%v", flowInfo.Id, jobId), string(jsonData[:]))
			if errPut != nil {
				errMsg = errPut.Error()
				return errPut
			}

			//sub result key
			jobResults = make(map[string]*common.JobResultInfo)
			wg.Add(1)
			go func(jobId string) {
				defer wg.Done()
				//get result
				for {
					watchChan := etcd.GetDefaultEtcdService().WatchKey("crontab", etcd.PrefixEnumJob, fmt.Sprintf("result/%v/%v", flowInfo.Id, jobId))
					for watch := range watchChan {
						if watch.Canceled {
							return
						}
						for _, ev := range watch.Events {
							if ev.IsModify() || ev.IsCreate() {
								var result *common.JobResultInfo
								errJson := json.Unmarshal(ev.Kv.Value, &result)
								if errJson != nil {
									log.Error("json unmarshal error", zap.Error(errJson))
								}
								jobResults[jobId] = result
								_ = etcd.GetDefaultEtcdService().DeleteByKey(string(ev.Kv.Key[:]))
								return
							}
						}
					}
				}
			}(jobId)
		}

		wg.Wait()

		for key, _ := range jobResults {
			jobOutput[key] = jobResults[key].JobOutput

			if jobResults[key] != nil && len(jobResults[key].ErrorMsg) > 0 {
				errMsg = jobResults[key].ErrorMsg
				return errors.New(errMsg)
			}
		}

		time.Sleep(15 * time.Second)
	}
	return nil
}
