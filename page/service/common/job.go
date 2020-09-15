package common

type CrontabJobListInputModel struct {
	PaginationInputModel
	Id      string `json:"id"`
	Name    string `json:"name"`
	JobType int32  `json:"job_type"`
	Status  int32  `json:"status"`
}

type CrontabJobDetailInputModel struct {
	Id string `json:"id"`
}

type CrontabJobListOutputModel struct {
	PaginationOutputModel
	Data []*CrontabJobModel `json:"data"`
}

type CrontabJobModel struct {
	Id         string `gorm:"not null; column:id" json:"id"`
	Name       string `gorm:"not null; column:name" json:"name"`
	JobType    int32  `gorm:"not null; column:job_type" json:"job_type"`
	Info       string `gorm:"not null; column:info" json:"info"`
	Desc       string `gorm:"not null; column:desc" json:"desc"`
	Status     int32  `gorm:"not null; column:status" json:"status"`
	CreateTime int64  `gorm:"not null; column:create_time" json:"create_time"`
	UpdateTime int64  `gorm:"not null; column:update_time" json:"update_time"`
}

type JobInfoBaseModel struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	JobType int32  `json:"job_type"`
}

type JobInfoCommandModel struct {
	JobInfoBaseModel
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

type JobInfoEnvironCalMetModel struct {
	JobInfoBaseModel
	JobDir     string `json:"job_dir"`
	Shell      string `json:"shell"`
	Arg        string `json:"arg"`
	Exe        string `json:"exe"`
	LogFile    string `json:"log_file"`
	InitExeDir string `json:"init_exe_dir"`
	InitConfig string `json:"init_config"`
	CalWRFDir  string `json:"cal_wrf_dir"`
}

type JobInfoEnvironCalPostModel struct {
	JobInfoBaseModel
	JobDir         string  `json:"job_dir"`
	Shell          string  `json:"shell"`
	Arg            string  `json:"arg"`
	Rec            []int32 `json:"rec"`
	LogFile        string  `json:"log_file"`
	InitExeFile    string  `json:"init_exe_file"`
	ModelDir       string  `json:"model_dir"`
	InitConfigFile string  `json:"init_config_file"`
	CalSumDir      string  `json:"cal_sum_dir"`
	PostUtilDir    string  `json:"post_util_dir"`
}

type JobInfoEnvironCalPuffModel struct {
	JobInfoBaseModel
	JobDir     string   `json:"job_dir"`
	Shell      string   `json:"shell"`
	Arg        string   `json:"arg"`
	Exe        []string `json:"exe"`
	LogFiles   []string `json:"log_files"`
	InitExeDir string   `json:"init_exe_dir"`
	CalMETDir  string   `json:"cal_met_dir"`
	InitConfig string   `json:"init_config"`
}

type JobInfoEnvironCalSumModel struct {
	JobInfoBaseModel
	JobDir     string `json:"job_dir"`
	Shell      string `json:"shell"`
	Arg        string `json:"arg"`
	Exe        string `json:"exe"`
	LogFile    string `json:"log_file"`
	InitExeDir string `json:"init_exe_dir"`
	CalPuffDir string `json:"cal_puff_dir"`
	InitConfig string `json:"init_config"`
}

type JobInfoEnvironCalWrfModel struct {
	JobInfoBaseModel
	JobDir     string `json:"job_dir"`
	Shell      string `json:"shell"`
	Arg        string `json:"arg"`
	Exe        string `json:"exe"`
	LogFile    string `json:"log_file"`
	WRFOutDir  string `json:"wrf_out_dir"`
	InitExeDir string `json:"init_exe_dir"`
	InitConfig string `json:"init_config"`
}

type JobInfoEnvironDownloadModel struct {
	JobInfoBaseModel
	JobDir       string `json:"job_dir"`
	UCARUsername string `json:"ucar_username"`
	UCARPassword string `json:"ucar_password"`
}

type JobInfoEnvironPostUtilModel struct {
	JobInfoBaseModel
	JobDir     string `json:"job_dir"`
	Shell      string `json:"shell"`
	Arg        string `json:"arg"`
	Exe        string `json:"exe"`
	LogFile    string `json:"log_file"`
	InitExeDir string `json:"init_exe_dir"`
	InitConfig string `json:"init_config"`
	CalSumDir  string `json:"cal_sum_dir"`
}

type JobInfoEnvironWpsWrfModel struct {
	JobInfoBaseModel
	JobDir            string              `json:"job_dir"`
	Shell             string              `json:"shell"`
	Arg               string              `json:"arg"`
	Command           string              `json:"command"`
	LogFiles          map[string][]string `json:"log_files"`
	WorkDir           string              `json:"work_dir"`
	InitConfigCommand string              `json:"init_config_command"`
	InitConfigWPS     []string            `json:"init_config_wps"`
	InitConfigWRF     []string            `json:"init_config_wrf"`
}

type JobInfoEnvironScrapyModel struct {
	JobInfoBaseModel
	JobDir string `json:"job_dir"`
}
