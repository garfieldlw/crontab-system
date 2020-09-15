package job

import (
	"github.com/garfieldlw/crontab-system/library/log"
	"github.com/garfieldlw/crontab-system/library/pgsql"
	"context"
	"errors"
	"go.uber.org/zap"
	"time"
)

type CrontabJobEntity struct {
	Id         string `gorm:"not null; column:id" json:"id"`
	Name       string `gorm:"not null; column:name" json:"name"`
	JobType    int32  `gorm:"not null; column:job_type" json:"job_type"`
	Info       string `gorm:"not null; column:info" json:"info"`
	Desc       string `gorm:"not null; column:desc" json:"desc"`
	Status     int32  `gorm:"not null; column:status" json:"status"`
	CreateTime int64  `gorm:"not null; column:create_time" json:"create_time"`
	UpdateTime int64  `gorm:"not null; column:update_time" json:"update_time"`
}

func (CrontabJobEntity) TableName() string {
	return "crontab_job"
}

func Create(ctx context.Context, id string, jobType int32, name, info, desc string, status int32) (*CrontabJobEntity, error) {
	conn := pgsql.GetDb()
	if conn == nil {
		return nil, errors.New("connect db fail")
	}

	in := new(CrontabJobEntity)
	in.Id = id
	in.JobType = jobType
	in.Name = name
	in.Info = info
	in.Desc = desc
	in.Status = status

	now := time.Now().Unix()
	in.CreateTime = now
	in.UpdateTime = now

	res := new(CrontabJobEntity)
	result := conn.Create(in).Scan(res)
	if result.Error != nil {
		log.Warn("process db fail", zap.Error(result.Error))
		return nil, result.Error
	}

	return res, nil
}

func Update(ctx context.Context, id string, jobType int32, name, info, desc string, status int32) (*CrontabJobEntity, error) {
	conn := pgsql.GetDb()
	if conn == nil {
		return nil, errors.New("connect db fail")
	}

	in := new(CrontabJobEntity)
	in.JobType = jobType
	in.Name = name
	in.Info = info
	in.Desc = desc
	in.Status = status

	now := time.Now().Unix()
	in.UpdateTime = now

	res := new(CrontabJobEntity)
	result := conn.Model(&CrontabJobEntity{}).Where(&CrontabJobEntity{Id: id}).Update(in).Scan(res)
	if result.Error != nil {
		log.Warn("process db fail", zap.Error(result.Error))
		return nil, result.Error
	}

	return res, nil
}

func UpdateStatus(ctx context.Context, id string, status int32) (*CrontabJobEntity, error) {
	conn := pgsql.GetDb()
	if conn == nil {
		return nil, errors.New("connect db fail")
	}

	if len(id) == 0 {
		return nil, errors.New("id is empty")
	}

	in := new(CrontabJobEntity)
	in.Status = status
	in.UpdateTime = time.Now().Unix()

	res := new(CrontabJobEntity)
	result := conn.Model(&CrontabJobEntity{}).Where(&CrontabJobEntity{Id: id}).Update(in).Scan(res)

	if result.Error != nil {
		log.Warn("process db fail", zap.Error(result.Error))
		return nil, result.Error
	}

	return res, nil
}

func DetailById(ctx context.Context, id string) (*CrontabJobEntity, error) {
	conn := pgsql.GetDb()
	if conn == nil {
		return nil, errors.New("connect db fail")
	}

	if len(id) == 0 {
		return nil, errors.New("id is empty")
	}

	res := new(CrontabJobEntity)
	result := conn.Model(&CrontabJobEntity{}).Where(&CrontabJobEntity{Id: id}).First(res)
	if result.Error != nil {
		log.Warn("process db fail", zap.Error(result.Error))
		return nil, result.Error
	}

	return res, nil
}

func DetailByIds(ctx context.Context, ids []string) ([]*CrontabJobEntity, error) {
	conn := pgsql.GetDb()
	if conn == nil {
		return nil, errors.New("connect db fail")
	}

	if len(ids) == 0 {
		return nil, nil
	}

	res := new([]*CrontabJobEntity)

	result := conn.Model(&CrontabJobEntity{}).Where("id in (?) ", ids).Find(&res)
	if result.Error != nil {
		log.Warn("process db fail", zap.Error(result.Error))
		return nil, result.Error
	}

	return *res, nil
}

func List(ctx context.Context, id string, jobType int32, name string, status int32, sortValue string, offset, limit int64) ([]*CrontabJobEntity, int64, error) {
	conn := pgsql.GetDb()
	if conn == nil {
		return nil, 0, errors.New("connect db fail")
	}

	if len(sortValue) == 0 {
		return nil, 0, errors.New("sort is empty")
	}

	condition := new(CrontabJobEntity)
	condition.Id = id
	condition.JobType = jobType
	condition.Name = name
	condition.Status = status

	query := conn.Model(&CrontabJobEntity{}).Where(condition)

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

	var listChan = make(chan []*CrontabJobEntity, 1)
	go func() {
		defer close(listChan)
		var actLimit []*CrontabJobEntity
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
	var actLimit []*CrontabJobEntity

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
