package common

import (
	"github.com/garfieldlw/crontab-system/page/service/common"
	"runtime"
	"time"
)

var (
	JobList []string
	JobChan chan *common.JobInfo
)

func init() {
	JobChan = make(chan *common.JobInfo, 10)
}

func GetEnvironJobDate(unix interface{}) time.Time {
	date := time.Unix(int64(unix.(float64)), 0).UTC()
	for {
		h := date.Hour()
		if h == 0 || h == 6 || h == 12 || h == 18 {
			date = date.Add(-time.Hour * 24)
			return date
		}

		date = date.Add(time.Hour)
	}
}

func GetFlowDir(flowInfo *common.FlowInfoEnvironModel) string {
	if runtime.GOOS == "windows" {
		return flowInfo.FlowDirWindows
	}

	return flowInfo.FlowDirLinux
}

func GetFlowInitDir(flowInfo *common.FlowInfoEnvironModel) string {
	if runtime.GOOS == "windows" {
		return flowInfo.InitDirWindows
	}

	return flowInfo.InitDirLinux
}
