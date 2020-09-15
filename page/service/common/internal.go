package common

type FlowInfo struct {
	FlowId  string             `json:"flow_id"`
	Force   bool               `json:"force"`
	DoForce map[int32]struct{} `json:"do_force"`
	Date    int64              `json:"date"`
}

type JobInfo struct {
	DoForce   bool                   `json:"do_force"`
	TraceId   int64                  `json:"trace_id"`
	FlowId    string                 `json:"flow_id"`
	FlowName  string                 `json:"flow_name"`
	FlowInfo  string                 `json:"flow_info"`
	FlowTye   int32                  `json:"flow_tye"`
	JobId     string                 `json:"job_id"`
	JobName   string                 `json:"job_name"`
	JobType   int32                  `json:"job_type"`
	JobInfo   string                 `json:"job_info"`
	FlowInput map[string]interface{} `json:"flow_input"`
	JobInput  map[string]interface{} `json:"job_input"`
}

type JobResultInfo struct {
	DoForce    bool                   `json:"do_force"`
	TraceId    int64                  `json:"trace_id"`
	FlowId     string                 `json:"flow_id"`
	FlowName   string                 `json:"flow_name"`
	FlowInfo   string                 `json:"flow_info"`
	FlowTye    int32                  `json:"flow_tye"`
	JobId      string                 `json:"job_id"`
	JobName    string                 `json:"job_name"`
	JobInfo    string                 `json:"job_info"`
	JobType    int32                  `json:"job_type"`
	StartTime  int64                  `json:"start_time"`
	EndTime    int64                  `json:"end_time"`
	FlowInput  map[string]interface{} `json:"flow_input"`
	JobInput   map[string]interface{} `json:"job_input"`
	FlowOutput map[string]interface{} `json:"flow_output"`
	JobOutput  map[string]interface{} `json:"job_output"`
	ErrorMsg   string                 `json:"error_msg"`
	Desc       string                 `json:"desc"`
}
