package log

import (
	"github.com/garfieldlw/crontab-system/library/unique"
	"github.com/garfieldlw/crontab-system/library/utils"
	"github.com/garfieldlw/crontab-system/page/model/log"
	"github.com/garfieldlw/crontab-system/page/model/unique-id"
	"github.com/garfieldlw/crontab-system/page/service/common"
	"context"
	"github.com/jinzhu/copier"
)

func CreateLogJob(ctx context.Context, id, fatherId, traceId int64, flowId, flowName, flowInfo string, jobId, jobName, jobInfo string, startTime, endTime int64, input, output, errMsg, desc string) (*common.CrontabLogJobModel, error) {
	ip := utils.LocalIP()

	if id == 0 {
		newId, errId := unique_id.GetUniqueId(unique.BizTypeLogJob)
		if errId != nil {
			return nil, errId
		}

		id = newId
	}

	if fatherId == 0 {
		fatherId = id
	}

	entity, errEntity := log.CreateLogJob(ctx, id, fatherId, traceId, ip, flowId, flowName, flowInfo, jobId, jobName, jobInfo, startTime, endTime, input, output, errMsg, desc)
	if errEntity != nil {
		return nil, errEntity
	}

	result := new(common.CrontabLogJobModel)
	_ = copier.Copy(result, entity)
	return result, nil
}

func UpdateLogJob(ctx context.Context, id, fatherId, traceId int64, flowId, flowName, flowInfo string, jobId, jobName, jobInfo string, startTime, endTime int64, input, output, errMsg, desc string) (*common.CrontabLogJobModel, error) {
	ip := utils.LocalIP()
	entity, errEntity := log.UpdateLogJob(ctx, id, fatherId, traceId, ip, flowId, flowName, flowInfo, jobId, jobName, jobInfo, startTime, endTime, input, output, errMsg, desc)
	if errEntity != nil {
		return nil, errEntity
	}

	result := new(common.CrontabLogJobModel)
	_ = copier.Copy(result, entity)
	return result, nil
}

func ListLogJob(ctx context.Context, id, fatherId, traceId int64, workIp, flowId, jobId, flowName, jobName string, startDate, endDate int64, sortValue string, offset, limit int64) (*common.CrontabLogJobListOutputModel, error) {
	entities, num, errEntities := log.ListLogJob(ctx, id, fatherId, traceId, workIp, flowId, jobId, flowName, jobName, startDate, endDate, sortValue, offset, limit)
	if errEntities != nil {
		return nil, errEntities
	}

	result := new(common.CrontabLogJobListOutputModel)
	result.Limit = limit
	result.Offset = offset
	result.Total = num

	for _, entity := range entities {
		data := new(common.CrontabLogJobModel)
		_ = copier.Copy(data, entity)
		result.Data = append(result.Data, data)
	}

	return result, nil
}

func CreateLogFlow(ctx context.Context, id, fatherId int64, flowId, flowName, flowInfo string, startTime, endTime int64, input, output, errMsg, desc string) (*common.CrontabLogFlowModel, error) {
	ip := utils.LocalIP()

	if id == 0 {
		newId, errId := unique_id.GetUniqueId(unique.BizTypeLogFlow)
		if errId != nil {
			return nil, errId
		}

		id = newId
	}

	if fatherId == 0 {
		fatherId = id
	}

	entity, errEntity := log.CreateLogFlow(ctx, id, fatherId, ip, flowId, flowName, flowInfo, startTime, endTime, input, output, errMsg, desc)
	if errEntity != nil {
		return nil, errEntity
	}

	result := new(common.CrontabLogFlowModel)
	_ = copier.Copy(result, entity)
	return result, nil
}

func UpdateLogFlow(ctx context.Context, id, fatherId int64, flowId, flowName, flowInfo string, startTime, endTime int64, input, output, errMsg, desc string) (*common.CrontabLogFlowModel, error) {
	ip := utils.LocalIP()
	entity, errEntity := log.UpdateLogFlow(ctx, id, fatherId, ip, flowId, flowName, flowInfo, startTime, endTime, input, output, errMsg, desc)
	if errEntity != nil {
		return nil, errEntity
	}

	result := new(common.CrontabLogFlowModel)
	_ = copier.Copy(result, entity)
	return result, nil
}

func ListLogFlow(ctx context.Context, id, fatherId int64, workIp, flowId, flowName string, startDate, endDate int64, sortValue string, offset, limit int64) (*common.CrontabLogFlowListOutputModel, error) {
	entities, num, errEntities := log.ListLogFlow(ctx, id, fatherId, workIp, flowId, flowName, startDate, endDate, sortValue, offset, limit)
	if errEntities != nil {
		return nil, errEntities
	}

	result := new(common.CrontabLogFlowListOutputModel)
	result.Limit = limit
	result.Offset = offset
	result.Total = num

	for _, entity := range entities {
		if entity == nil || entity.Id == 0 {
			continue
		}

		data := new(common.CrontabLogFlowModel)
		_ = copier.Copy(data, entity)

		logJobs, _, _ := log.ListLogJob(ctx, 0, 0, entity.Id, "", "", "", "", "", 0, 0, "id desc", 0, 1000)

		for _, logJob := range logJobs {
			if logJob == nil || logJob.Id == 0 {
				continue
			}

			l := new(common.CrontabLogJobModel)
			_ = copier.Copy(l, logJob)

			data.LogJob = append(data.LogJob, l)
		}

		result.Data = append(result.Data, data)
	}

	return result, nil
}
