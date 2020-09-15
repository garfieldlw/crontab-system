package handler

import (
	"context"
	"errors"
	"github.com/garfieldlw/crontab-system/cron/worker/service/handler/command"
	"github.com/garfieldlw/crontab-system/library/log"
	"github.com/garfieldlw/crontab-system/page/service/common"
	"go.uber.org/zap"
)

var handlers = make(map[string]DoJobInterface)

func Do(ctx context.Context, info *common.JobInfo) (map[string]interface{}, error) {
	if len(info.JobId) == 0 {
		return nil, errors.New("invalid group")
	}

	key := ""
	switch info.JobType {
	case 1:
		{
			return nil, errors.New("it is delete action")
		}
	case 2: //command
		key = "Command"
	default:
		key = info.JobId
	}

	doFunc, ok := handlers[key]
	if !ok {
		log.Warn("process handler, load function fail", zap.Any("job info", info))
		return nil, errors.New("load function fail")
	}

	if doFunc == nil {
		return nil, errors.New("load function is empty")
	}

	if !info.DoForce {
		result, doCheck, errDoCheck := doFunc.CheckDo(ctx, info.JobId, info.JobType, info.JobName, info.FlowInfo, info.JobInfo, info.FlowInput, info.JobInput)
		if errDoCheck != nil {
			log.Warn("process handler, do check error", zap.Error(errDoCheck))
			return nil, errDoCheck
		}

		if doCheck {
			return result, nil
		}
	}

	result, errResult := doFunc.Do(ctx, info.JobId, info.JobType, info.JobName, info.FlowInfo, info.JobInfo, info.FlowInput, info.JobInput)
	if errResult != nil {
		log.Warn("process handler, call function error", zap.Error(errResult))
		return nil, errResult
	}

	return result, nil
}

func init() {
	handlers["Command"] = new(command.Handler)
}
