package common

import "github.com/garfieldlw/crontab-system/page/service/common"

var (
	FlowChan chan *common.FlowInfo
)

func init() {
	FlowChan = make(chan *common.FlowInfo, 10)
}
