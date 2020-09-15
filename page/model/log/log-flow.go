package log

import (
	"github.com/garfieldlw/crontab-system/library/log"
	"github.com/garfieldlw/crontab-system/library/pgsql"
	"context"
	"errors"
	"go.uber.org/zap"
	"time"
)

type CrontabLogFlowEntity struct {
	Id         int64  `gorm:"not null; column:id" json:"id"`
	FatherId   int64  `gorm:"not null; column:father_id" json:"father_id"`
	WorkerIp   string `gorm:"not null; column:worker_ip" json:"worker_ip"`
	FlowId     string `gorm:"not null; column:flow_id" json:"flow_id"`
	FlowName   string `gorm:"not null; column:flow_name" json:"flow_name"`
	FlowInfo   string `gorm:"not null; column:flow_info" json:"flow_info"`
	StartTime  int64  `gorm:"not null; column:start_time" json:"start_time"`
	EndTime    int64  `gorm:"not null; column:end_time" json:"end_time"`
	Input      string `gorm:"not null; column:input" json:"input"`
	Output     string `gorm:"not null; column:output" json:"output"`
	ErrorMsg   string `gorm:"not null; column:error_msg" json:"error_msg"`
	Desc       string `gorm:"not null; column:desc" json:"desc"`
	CreateTime int64  `gorm:"not null; column:create_time" json:"create_time"`
	UpdateTime int64  `gorm:"not null; column:update_time" json:"update_time"`
}

func (CrontabLogFlowEntity) TableName() string {
	return "crontab_log_flow"
}

func CreateLogFlow(ctx context.Context, id, fatherId int64, workerIp string, flowId, flowName, flowInfo string, startTime, endTime int64, input, output, errMsg, desc string) (*CrontabLogFlowEntity, error) {
	conn := pgsql.GetDb()
	if conn == nil {
		return nil, errors.New("connect db fail")
	}

	in := new(CrontabLogFlowEntity)
	in.Id = id
	in.FatherId = fatherId
	in.WorkerIp = workerIp
	in.FlowId = flowId
	in.FlowName = flowName
	in.FlowInfo = flowInfo
	in.StartTime = startTime
	in.EndTime = endTime
	in.Input = input
	in.Output = output
	in.ErrorMsg = errMsg
	in.Desc = desc

	now := time.Now().Unix()
	in.CreateTime = now
	in.UpdateTime = now

	res := new(CrontabLogFlowEntity)
	result := conn.Create(in).Scan(res)
	if result.Error != nil {
		log.Warn("process db fail", zap.Error(result.Error))
		return nil, result.Error
	}

	return res, nil
}

func UpdateLogFlow(ctx context.Context, id, fatherId int64, workerIp string, flowId, flowName, flowInfo string, startTime, endTime int64, input, output, errMsg, desc string) (*CrontabLogFlowEntity, error) {
	conn := pgsql.GetDb()
	if conn == nil {
		return nil, errors.New("connect db fail")
	}

	in := new(CrontabLogFlowEntity)
	in.FatherId = fatherId
	in.WorkerIp = workerIp
	in.FlowId = flowId
	in.FlowName = flowName
	in.FlowInfo = flowInfo
	in.StartTime = startTime
	in.EndTime = endTime
	in.Input = input
	in.Output = output
	in.ErrorMsg = errMsg
	in.Desc = desc
	in.UpdateTime = time.Now().Unix()

	res := new(CrontabLogFlowEntity)
	result := conn.Model(&CrontabLogFlowEntity{}).Where(&CrontabLogFlowEntity{Id: id}).Update(in).Scan(res)
	if result.Error != nil {
		log.Warn("process db fail", zap.Error(result.Error))
		return nil, result.Error
	}

	return res, nil
}

func ListLogFlow(ctx context.Context, id, fatherId int64, workIp, flowId, flowName string, startTime, endTime int64, sortValue string, offset, limit int64) ([]*CrontabLogFlowEntity, int64, error) {
	conn := pgsql.GetDb()
	if conn == nil {
		return nil, 0, errors.New("connect db fail")
	}

	if len(sortValue) == 0 {
		return nil, 0, errors.New("sort is empty")
	}

	condition := new(CrontabLogFlowEntity)
	condition.Id = id
	condition.FatherId = fatherId
	condition.FlowId = flowId
	condition.WorkerIp = workIp
	condition.FlowName = flowName

	query := conn.Model(&CrontabLogFlowEntity{}).Where(condition)

	if startTime > 0 {
		query = query.Where(" \"start_time\" >= ?", startTime)
	}

	if endTime > 0 {
		query = query.Where(" \"start_time\" < ?", endTime)
	}

	var countChan = make(chan int64, 1)
	go func() {
		defer close(countChan)
		var num int64
		resCount := query.Count(&num)
		if resCount.Error != nil {
			log.Warn("query error", zap.Error(resCount.Error))
			countChan <- 0
			return
		}
		countChan <- num
		return
	}()

	var listChan = make(chan []*CrontabLogFlowEntity, 1)
	go func() {
		defer close(listChan)
		var actLimit []*CrontabLogFlowEntity
		res := query.Offset(offset).Limit(limit).Order(sortValue).Find(&actLimit)
		if res.Error != nil {
			log.Warn("query error", zap.Error(res.Error))
			listChan <- nil
			return
		}
		listChan <- actLimit
		return
	}()

	var num int64
	var actLimit []*CrontabLogFlowEntity

	select {
	case num = <-countChan:
	case <-time.After(time.Second * 10):
		return nil, 0, errors.New("db query timeout")
	}

	select {
	case actLimit = <-listChan:
	case <-time.After(time.Second * 10):
		return nil, 0, errors.New("db query timeout")
	}

	return actLimit, num, nil
}
