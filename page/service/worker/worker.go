package worker

import (
	"github.com/garfieldlw/crontab-system/page/model/worker"
	"github.com/garfieldlw/crontab-system/page/service/common"
	"context"
	"github.com/jinzhu/copier"
)

func DetailById(ctx context.Context, id int64) (*common.CrontabWorkerModel, error) {
	entity, errEntity := worker.DetailById(ctx, id)
	if errEntity != nil {
		return nil, errEntity
	}

	result := new(common.CrontabWorkerModel)
	_ = copier.Copy(result, entity)
	return result, nil
}

func DetailByNameIpOs(ctx context.Context, name, ip, os string) (*common.CrontabWorkerModel, error) {
	entity, errEntity := worker.DetailByKeys(ctx, name, ip, os)
	if errEntity != nil {
		return nil, errEntity
	}

	result := new(common.CrontabWorkerModel)
	_ = copier.Copy(result, entity)
	return result, nil
}
