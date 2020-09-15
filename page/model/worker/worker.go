package worker

import (
	"github.com/garfieldlw/crontab-system/library/log"
	"github.com/garfieldlw/crontab-system/library/pgsql"
	"context"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"go.uber.org/zap"
	"time"
)

type CrontabWorkerEntity struct {
	Id           int64          `gorm:"not null; column:id" json:"id"`
	Name         string         `gorm:"not null; column:name" json:"name"`
	OS           string         `gorm:"not null; column:os" json:"os"`
	ARCH         string         `gorm:"not null; column:arch" json:"arch"`
	Ip           string         `gorm:"not null; column:ip" json:"ip"`
	JobList      pq.StringArray `gorm:"type:varchar[]; not null; column:job_list" json:"job_list"`
	StatusOnline int32          `gorm:"not null; column:status_online" json:"status_online"`
	Status       int32          `gorm:"not null; column:status" json:"status"`
	CreateTime   int64          `gorm:"not null; column:create_time" json:"create_time"`
	UpdateTime   int64          `gorm:"not null; column:update_time" json:"update_time"`
}

func (CrontabWorkerEntity) TableName() string {
	return "crontab_worker"
}

func Upset(ctx context.Context, id int64, name, ip, os, arch string, jobList []string, statusOnline, status int32) (*CrontabWorkerEntity, error) {
	conn := pgsql.GetDb()
	if conn == nil {
		return nil, errors.New("connect db fail")
	}

	now := time.Now().Unix()
	sql := fmt.Sprintf("INSERT INTO \"crontab_worker\" (id, name, ip, os, arch, job_list, status_online, status, create_time, update_time) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) ON CONFLICT (name,ip,os) DO UPDATE SET update_time = %v;", now)

	res := new(CrontabWorkerEntity)
	result := conn.Exec(sql, id, name, ip, os, arch, jobList, statusOnline, status, now, now).First(res)
	if result.Error != nil {
		log.Warn("process db fail", zap.Error(result.Error))
		return nil, result.Error
	}

	updateValue := new(CrontabWorkerEntity)
	updateValue.Name = name
	updateValue.Ip = ip
	updateValue.OS = os
	updateValue.JobList = jobList
	updateValue.StatusOnline = statusOnline
	updateValue.Status = status
	updateValue.UpdateTime = now

	resultUpdate := conn.Model(&CrontabWorkerEntity{}).Where(&CrontabWorkerEntity{Id: id}).Update(updateValue).Scan(res)
	if resultUpdate.Error != nil {
		log.Warn("process db fail", zap.Error(result.Error))
		return nil, resultUpdate.Error
	}

	return res, nil
}

func Create(ctx context.Context, id int64, name, ip, os, arch string, jobList []string, statusOnline, status int32) (*CrontabWorkerEntity, error) {
	conn := pgsql.GetDb()
	if conn == nil {
		return nil, errors.New("connect db fail")
	}

	in := new(CrontabWorkerEntity)
	in.Id = id
	in.Name = name
	in.Ip = ip
	in.OS = os
	in.ARCH = arch
	in.JobList = jobList
	in.StatusOnline = statusOnline
	in.Status = status

	now := time.Now().Unix()
	in.CreateTime = now
	in.UpdateTime = now

	res := new(CrontabWorkerEntity)
	result := conn.Create(in).Scan(res)
	if result.Error != nil {
		log.Warn("process db fail", zap.Error(result.Error))
		return nil, result.Error
	}

	return res, nil
}

func Update(ctx context.Context, id int64, name, ip, os, arch string, jobList []string, statusOnline, status int32) (*CrontabWorkerEntity, error) {
	conn := pgsql.GetDb()
	if conn == nil {
		return nil, errors.New("connect db fail")
	}

	in := new(CrontabWorkerEntity)
	in.Name = name
	in.Ip = ip
	in.OS = os
	in.ARCH = arch
	in.JobList = jobList
	in.StatusOnline = statusOnline
	in.Status = status

	now := time.Now().Unix()
	in.UpdateTime = now

	res := new(CrontabWorkerEntity)
	result := conn.Model(&CrontabWorkerEntity{}).Where(&CrontabWorkerEntity{Id: id}).Update(in).Scan(res)
	if result.Error != nil {
		log.Warn("process db fail", zap.Error(result.Error))
		return nil, result.Error
	}

	return res, nil
}

func UpdateStatus(ctx context.Context, id int64, status int32) (*CrontabWorkerEntity, error) {
	conn := pgsql.GetDb()
	if conn == nil {
		return nil, errors.New("connect db fail")
	}

	if id == 0 {
		return nil, errors.New("id is 0")
	}

	in := new(CrontabWorkerEntity)
	in.Status = status
	in.UpdateTime = time.Now().Unix()

	res := new(CrontabWorkerEntity)
	result := conn.Model(&CrontabWorkerEntity{}).Where(&CrontabWorkerEntity{Id: id}).Update(in).Scan(res)

	if result.Error != nil {
		log.Warn("process db fail", zap.Error(result.Error))
		return nil, result.Error
	}

	return res, nil
}

func DetailById(ctx context.Context, id int64) (*CrontabWorkerEntity, error) {
	conn := pgsql.GetDb()
	if conn == nil {
		return nil, errors.New("connect db fail")
	}

	if id == 0 {
		return nil, errors.New("id is 0")
	}

	res := new(CrontabWorkerEntity)
	result := conn.Model(&CrontabWorkerEntity{}).First(res, id)
	if result.Error != nil {
		log.Warn("process db fail", zap.Error(result.Error))
		return nil, result.Error
	}

	return res, nil
}

func DetailByKeys(ctx context.Context, name, ip, os string) (*CrontabWorkerEntity, error) {
	conn := pgsql.GetDb()
	if conn == nil {
		return nil, errors.New("connect db fail")
	}

	res := new(CrontabWorkerEntity)
	result := conn.Model(&CrontabWorkerEntity{}).Where(&CrontabWorkerEntity{Name: name, Ip: ip, OS: os}).First(res)
	if result.Error != nil {
		log.Warn("process db fail", zap.Error(result.Error))
		return nil, result.Error
	}

	return res, nil
}

func DetailByIds(ctx context.Context, ids []int64) ([]*CrontabWorkerEntity, error) {
	conn := pgsql.GetDb()
	if conn == nil {
		return nil, errors.New("connect db fail")
	}

	if len(ids) == 0 {
		return nil, nil
	}

	res := new([]*CrontabWorkerEntity)

	result := conn.Model(&CrontabWorkerEntity{}).Where("id in (?) ", ids).Find(&res)
	if result.Error != nil {
		log.Warn("process db fail", zap.Error(result.Error))
		return nil, result.Error
	}

	return *res, nil
}

func List(ctx context.Context, id int64, name, ip, os, arch string, statusOnline, status int32, sortValue string, offset, limit int64) ([]*CrontabWorkerEntity, int64, error) {
	conn := pgsql.GetDb()
	if conn == nil {
		return nil, 0, errors.New("connect db fail")
	}

	if len(sortValue) == 0 {
		return nil, 0, errors.New("sort is empty")
	}

	condition := new(CrontabWorkerEntity)
	condition.Id = id
	condition.Name = name
	condition.Ip = ip
	condition.OS = os
	condition.ARCH = arch
	condition.StatusOnline = statusOnline
	condition.Status = status

	query := conn.Model(&CrontabWorkerEntity{}).Where(condition)

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

	var listChan = make(chan []*CrontabWorkerEntity, 1)
	go func() {
		defer close(listChan)
		var actLimit []*CrontabWorkerEntity
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
	var actLimit []*CrontabWorkerEntity

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
