package common

type CrontabFlowListInputModel struct {
	PaginationInputModel
	Id       string `json:"id"`
	Name     string `json:"name"`
	FlowType int32  `json:"flow_type"`
	Status   int32  `json:"status"`
}

type CrontabFlowDetailInputModel struct {
	Id string `json:"id"`
}

type CrontabFlowListOutputModel struct {
	PaginationOutputModel
	Data []*CrontabFlowModel `json:"data"`
}

type CrontabFlowModel struct {
	Id         string `gorm:"not null; column:id" json:"id"`
	Name       string `gorm:"not null; column:name" json:"name"`
	Info       string `gorm:"not null; column:info" json:"info"`
	FlowType   int32  `gorm:"not null; column:flow_type" json:"flow_type"`
	Spec       string `gorm:"not null; column:spec" json:"spec"`
	Desc       string `gorm:"not null; column:desc" json:"desc"`
	Status     int32  `gorm:"not null; column:status" json:"status"`
	CreateTime int64  `gorm:"not null; column:create_time" json:"create_time"`
	UpdateTime int64  `gorm:"not null; column:update_time" json:"update_time"`
}

type FlowInfoModel struct {
	Jobs map[int32][]string `json:"jobs"`
}

type FlowInfoEnvironModel struct {
	FlowInfoModel
	FlowDirLinux   string `json:"flow_dir_linux"`
	InitDirLinux   string `json:"init_dir_linux"`
	FlowDirWindows string `json:"flow_dir_windows"`
	InitDirWindows string `json:"init_dir_windows"`
}

type FlowInfoEnvironScrapyModel struct {
	FlowInfoModel
}

type FlowDoInputModel struct {
	FlowId  string             `json:"flow_id"`
	DoForce map[int32]struct{} `json:"do_force"`
	Date    int64              `json:"date"`
}
