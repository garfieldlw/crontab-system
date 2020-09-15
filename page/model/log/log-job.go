package log

import (
	"github.com/garfieldlw/crontab-system/library/log"
	"github.com/garfieldlw/crontab-system/library/pgsql"
	"context"
	"errors"
	"go.uber.org/zap"
	"time"
)

type CrontabLogJobEntity struct {
	Id         int64  `gorm:"not null; column:id" json:"id"`
	FatherId   int64  `gorm:"not null; column:father_id" json:"father_id"`
	TraceId    int64  `gorm:"not null; column:trace_id" json:"trace_id"`
	WorkerIp   string `gorm:"not null; column:worker_ip" json:"worker_ip"`
	FlowId     string `gorm:"not null; column:flow_id" json:"flow_id"`
	FlowName   string `gorm:"not null; column:flow_name" json:"flow_name"`
	FlowInfo   string `gorm:"not null; column:flow_info" json:"flow_info"`
	JobId      string `gorm:"not null; column:job_id" json:"job_id"`
	JobName    string `gorm:"not null; column:job_name" json:"job_name"`
	JobInfo    string `gorm:"not null; column:job_info" json:"job_info"`
	StartTime  int64  `gorm:"not null; column:start_time" json:"start_time"`
	EndTime    int64  `gorm:"not null; column:end_time" json:"end_time"`
	Input      string `gorm:"not null; column:input" json:"input"`
	Output     string `gorm:"not null; column:output" json:"output"`
	ErrorMsg   string `gorm:"not null; column:error_msg" json:"error_msg"`
	Desc       string `gorm:"not null; column:desc" json:"desc"`
	CreateTime int64  `gorm:"not null; column:create_time" json:"create_time"`
	UpdateTime int64  `gorm:"not null; column:update_time" json:"update_time"`
}

func (CrontabLogJobEntity) TableName() string {
	return "crontab_log_job"
}

func CreateLogJob(ctx context.Context, id, fatherId, traceId int64, workerIp string, flowId, flowName, flowInfo string, jobId, jobName, jobInfo string, startTime, endTime int64, input, output, errMsg, desc string) (*CrontabLogJobEntity, error) {
	conn := pgsql.GetDb()
	if conn == nil {
		return nil, errors.New("connect db fail")
	}

	in := new(CrontabLogJobEntity)
	in.Id = id
	in.FatherId = fatherId
	in.TraceId = traceId
	in.WorkerIp = workerIp
	in.FlowId = flowId
	in.FlowName = flowName
	in.FlowInfo = flowInfo
	in.JobId = jobId
	in.JobName = jobName
	in.JobInfo = jobInfo
	in.StartTime = startTime
	in.EndTime = endTime
	in.Input = input
	in.Output = output
	in.ErrorMsg = errMsg
	in.Desc = desc

	now := time.Now().Unix()
	in.CreateTime = now
	in.UpdateTime = now

	res := new(CrontabLogJobEntity)
	result := conn.Create(in).Scan(res)
	if result.Error != nil {
		log.Warn("process db fail", zap.Error(result.Error))
		return nil, result.Error
	}

	return res, nil
}

func UpdateLogJob(ctx context.Context, id, fatherId, traceId int64, workerIp string, flowId, flowName, flowInfo string, jobId, jobName, jobInfo string, startTime, endTime int64, input, output, errMsg, desc string) (*CrontabLogJobEntity, error) {
	conn := pgsql.GetDb()
	if conn == nil {
		return nil, errors.New("connect db fail")
	}

	in := new(CrontabLogJobEntity)
	in.FatherId = fatherId
	in.TraceId = traceId
	in.WorkerIp = workerIp
	in.FlowId = flowId
	in.FlowName = flowName
	in.FlowInfo = flowInfo
	in.JobId = jobId
	in.JobName = jobName
	in.JobInfo = jobInfo
	in.StartTime = startTime
	in.EndTime = endTime
	in.Input = input
	in.Output = output
	in.ErrorMsg = errMsg
	in.Desc = desc
	in.UpdateTime = time.Now().Unix()

	res := new(CrontabLogJobEntity)
	result := conn.Model(&CrontabLogJobEntity{}).Where(&CrontabLogJobEntity{Id: id}).Update(in).Scan(res)
	if result.Error != nil {
		log.Warn("process db fail", zap.Error(result.Error))
		return nil, result.Error

	}

	return res, nil
}

func ListLogJob(ctx context.Context, id, fatherId, traceId int64, workIp, flowId, jobId, flowName, jobName string, startTime, endTime int64, sortValue string, offset, limit int64) ([]*CrontabLogJobEntity, int64, error) {
	conn := pgsql.GetDb()
	if conn == nil {
		return nil, 0, errors.New("connect db fail")
	}

	if len(sortValue) == 0 {
		return nil, 0, errors.New("sort is empty")
	}

	condition := new(CrontabLogJobEntity)
	condition.Id = id
	condition.FatherId = fatherId
	condition.TraceId = traceId
	condition.FlowId = flowId
	condition.JobId = jobId
	condition.WorkerIp = workIp
	condition.FlowName = flowName
	condition.JobName = jobName

	query := conn.Model(&CrontabLogJobEntity{}).Where(condition)

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

	var listChan = make(chan []*CrontabLogJobEntity, 1)
	go func() {
		defer close(listChan)
		var actLimit []*CrontabLogJobEntity
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
	var actLimit []*CrontabLogJobEntity

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
