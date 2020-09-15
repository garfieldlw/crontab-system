package common

import "github.com/garfieldlw/crontab-system/library/bigint"

type CrontabLogJobListInputModel struct {
	PaginationInputModel
	Id        *bigint.BigInt `json:"id"`
	FatherId  *bigint.BigInt `json:"father_id"`
	TraceId   *bigint.BigInt `json:"trace_id"`
	WorkIp    string         `json:"name"`
	FlowId    string         `json:"flow_id"`
	FlowName  string         `json:"flow_name"`
	JobId     string         `json:"job_id"`
	JobName   string         `json:"job_name"`
	StartTime int64          `json:"start_time"`
	EndTime   int64          `json:"end_time"`
}

type CrontabLogJobListOutputModel struct {
	PaginationOutputModel
	Data []*CrontabLogJobModel `json:"data"`
}

type CrontabLogFlowListInputModel struct {
	PaginationInputModel
	Id        *bigint.BigInt `json:"id"`
	FatherId  *bigint.BigInt `json:"father_id"`
	WorkIp    string         `json:"name"`
	FlowId    string         `json:"flow_id"`
	FlowName  string         `json:"flow_name"`
	StartTime int64          `json:"start_time"`
	EndTime   int64          `json:"end_time"`
}

type CrontabLogFlowListOutputModel struct {
	PaginationOutputModel
	Data []*CrontabLogFlowModel `json:"data"`
}

type CrontabLogJobModel struct {
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

type CrontabLogFlowModel struct {
	Id         int64                 `gorm:"not null; column:id" json:"id"`
	FatherId   int64                 `gorm:"not null; column:father_id" json:"father_id"`
	WorkerIp   string                `gorm:"not null; column:worker_ip" json:"worker_ip"`
	FlowId     string                `gorm:"not null; column:flow_id" json:"flow_id"`
	FlowName   string                `gorm:"not null; column:flow_name" json:"flow_name"`
	FlowInfo   string                `gorm:"not null; column:flow_info" json:"flow_info"`
	StartTime  int64                 `gorm:"not null; column:start_time" json:"start_time"`
	EndTime    int64                 `gorm:"not null; column:end_time" json:"end_time"`
	Input      string                `gorm:"not null; column:input" json:"input"`
	Output     string                `gorm:"not null; column:output" json:"output"`
	ErrorMsg   string                `gorm:"not null; column:error_msg" json:"error_msg"`
	Desc       string                `gorm:"not null; column:desc" json:"desc"`
	CreateTime int64                 `gorm:"not null; column:create_time" json:"create_time"`
	UpdateTime int64                 `gorm:"not null; column:update_time" json:"update_time"`
	LogJob     []*CrontabLogJobModel `json:"log_job"`
}
