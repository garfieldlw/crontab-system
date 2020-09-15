package handler

import (
	"context"
)

type DoJobInterface interface {
	Name() string
	CheckDo(ctx context.Context, jobId string, jobType int32, jobName, flowInfo, jobInfo string, flowInput, jobInput map[string]interface{}) (map[string]interface{}, bool, error)
	Do(ctx context.Context, jobId string, jobType int32, jobName, flowInfo, jobInfo string, flowInput, jobInput map[string]interface{}) (map[string]interface{}, error)
}
