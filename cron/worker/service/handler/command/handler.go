package command

import (
	"github.com/garfieldlw/crontab-system/page/service/common"
	"context"
	"encoding/json"
	"os/exec"
)

type Handler struct {
}

func (handler *Handler) Name() string {
	return "Command"
}

func (handler *Handler) CheckDo(ctx context.Context, jobId string, jobType int32, jobName, flowInfo, jobInfo string, flowInput, jobInput map[string]interface{}) (map[string]interface{}, bool, error) {
	return nil, false, nil
}

func (handler *Handler) Do(ctx context.Context, jobId string, jobType int32, jobName, flowInfo, jobInfo string, flowInput, jobInput map[string]interface{}) (map[string]interface{}, error) {
	var info *common.JobInfoCommandModel
	errJson := json.Unmarshal([]byte(jobInfo[:]), &info)
	if errJson != nil {
		return nil, errJson
	}

	cmd := exec.CommandContext(ctx, info.Command, info.Args...)
	output, errOutput := cmd.CombinedOutput()
	if errOutput != nil {
		return nil, errOutput
	}

	result := map[string]interface{}{"": string(output[:])}
	return result, nil
}
