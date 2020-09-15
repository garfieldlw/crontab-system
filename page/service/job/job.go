package job

import (
	"github.com/garfieldlw/crontab-system/page/model/job"
	"github.com/garfieldlw/crontab-system/page/service/common"
	"context"
	"github.com/jinzhu/copier"
)

func Create(ctx context.Context, id string, jobType int32, name, info, desc string, status int32) (*common.CrontabJobModel, error) {
	entity, errEntity := job.Create(ctx, id, jobType, name, info, desc, status)
	if errEntity != nil {
		return nil, errEntity
	}

	res := new(common.CrontabJobModel)
	_ = copier.Copy(res, entity)
	return res, nil
}

func Update(ctx context.Context, id string, jobType int32, name, info, desc string, status int32) (*common.CrontabJobModel, error) {
	entity, errEntity := job.Update(ctx, id, jobType, name, info, desc, status)
	if errEntity != nil {
		return nil, errEntity
	}

	res := new(common.CrontabJobModel)
	_ = copier.Copy(res, entity)
	return res, nil
}

func UpdateStatus(ctx context.Context, id string, status int32) (*common.CrontabJobModel, error) {
	entity, errEntity := job.UpdateStatus(ctx, id, status)
	if errEntity != nil {
		return nil, errEntity
	}

	res := new(common.CrontabJobModel)
	_ = copier.Copy(res, entity)
	return res, nil
}

func DetailById(ctx context.Context, id string) (*common.CrontabJobModel, error) {
	entity, errEntity := job.DetailById(ctx, id)
	if errEntity != nil {
		return nil, errEntity
	}

	res := new(common.CrontabJobModel)
	_ = copier.Copy(res, entity)
	return res, nil
}

func List(ctx context.Context, id string, jobType int32, name string, status int32, sortValue string, offset, limit int64) (*common.CrontabJobListOutputModel, error) {
	entities, num, errEntities := job.List(ctx, id, jobType, name, status, sortValue, offset, limit)
	if errEntities != nil {
		return nil, errEntities
	}

	result := new(common.CrontabJobListOutputModel)
	result.Limit = limit
	result.Offset = offset
	result.Total = num

	for _, entity := range entities {
		data := new(common.CrontabJobModel)
		_ = copier.Copy(data, entity)
		result.Data = append(result.Data, data)
	}

	return result, nil
}
