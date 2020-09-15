package flow

import (
	"github.com/garfieldlw/crontab-system/page/model/flow"
	"github.com/garfieldlw/crontab-system/page/service/common"
	"context"
	"github.com/jinzhu/copier"
)

func Create(ctx context.Context, id, name, info, spec, desc string, flowType, status int32) (*common.CrontabFlowModel, error) {
	entity, errEntity := flow.Create(ctx, id, name, info, spec, desc, flowType, status)
	if errEntity != nil {
		return nil, errEntity
	}

	res := new(common.CrontabFlowModel)
	_ = copier.Copy(res, entity)
	return res, nil
}

func Update(ctx context.Context, id, name, info, spec, desc string, flowType, status int32) (*common.CrontabFlowModel, error) {
	entity, errEntity := flow.Update(ctx, id, name, info, spec, desc, flowType, status)
	if errEntity != nil {
		return nil, errEntity
	}

	res := new(common.CrontabFlowModel)
	_ = copier.Copy(res, entity)
	return res, nil
}

func UpdateStatus(ctx context.Context, id string, status int32) (*common.CrontabFlowModel, error) {
	entity, errEntity := flow.UpdateStatus(ctx, id, status)
	if errEntity != nil {
		return nil, errEntity
	}

	res := new(common.CrontabFlowModel)
	_ = copier.Copy(res, entity)
	return res, nil
}

func DetailById(ctx context.Context, id string) (*common.CrontabFlowModel, error) {
	entity, errEntity := flow.DetailById(ctx, id)
	if errEntity != nil {
		return nil, errEntity
	}

	res := new(common.CrontabFlowModel)
	_ = copier.Copy(res, entity)
	return res, nil
}

func List(ctx context.Context, id, name string, flowType, status int32, sortValue string, offset, limit int64) (*common.CrontabFlowListOutputModel, error) {
	entities, num, errEntities := flow.List(ctx, id, name, flowType, status, sortValue, offset, limit)
	if errEntities != nil {
		return nil, errEntities
	}

	result := new(common.CrontabFlowListOutputModel)
	result.Limit = limit
	result.Offset = offset
	result.Total = num

	for _, entity := range entities {
		data := new(common.CrontabFlowModel)
		_ = copier.Copy(data, entity)
		result.Data = append(result.Data, data)
	}

	return result, nil
}
