package flow

import (
	"github.com/garfieldlw/crontab-system/library/etcd"
	"github.com/garfieldlw/crontab-system/page/service/common"
	"context"
	"encoding/json"
)

func Do(ctx context.Context, flowId string, doForce map[int32]struct{}, date int64) error {
	val := new(common.FlowInfo)
	if flowId == "Environ" {
		val.Date = date
		val.Force = false
		val.DoForce = doForce
	}
	val.FlowId = flowId

	data, errJson := json.Marshal(val)
	if errJson != nil {
		return errJson
	}

	return etcd.GetDefaultEtcdService().Put("crontab", etcd.PrefixEnumFlow, "list/"+flowId, string(data[:]))
}
