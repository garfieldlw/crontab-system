package common

type CrontabWorkerModel struct {
	Id           int64    `gorm:"not null; column:id" json:"id"`
	Name         string   `gorm:"not null; column:name" json:"name"`
	OS           string   `gorm:"not null; column:os" json:"os"`
	ARCH         string   `gorm:"not null; column:arch" json:"arch"`
	Ip           string   `gorm:"not null; column:ip" json:"ip"`
	JobList      []string `gorm:"type:varchar[]; not null; column:job_list" json:"job_list"`
	StatusOnline int32    `gorm:"not null; column:status_online" json:"status_online"`
	Status       int32    `gorm:"not null; column:status" json:"status"`
	CreateTime   int64    `gorm:"not null; column:create_time" json:"create_time"`
	UpdateTime   int64    `gorm:"not null; column:update_time" json:"update_time"`
}
