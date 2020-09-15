package init_cron

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/garfieldlw/crontab-system/cron/worker/service/common"
	"github.com/garfieldlw/crontab-system/cron/worker/service/job"
	"github.com/garfieldlw/crontab-system/library/etcd"
	"github.com/garfieldlw/crontab-system/library/log"
	"github.com/garfieldlw/crontab-system/library/utils"
	common2 "github.com/garfieldlw/crontab-system/page/service/common"
	"github.com/garfieldlw/crontab-system/page/service/worker"
	"go.uber.org/zap"
	"time"
)

var workId int64 = 0

func Worker(ctx context.Context) {
	hostname := utils.Hostname()
	ip := utils.LocalIP()
	os, _ := utils.OSInfo()

	fmt.Println(hostname, ip, os)
	w, errWorker := worker.DetailByNameIpOs(ctx, hostname, ip, os)
	if errWorker != nil || w == nil || w.Id == 0 {
		log.Panic("worker config is empty", zap.Any("worker info", w), zap.Error(errWorker))
	}

	workId = w.Id

	common.JobList = w.JobList
	go func() {
		ticker := time.NewTicker(time.Second * 60 * 5)
		defer ticker.Stop()
		for {
			<-ticker.C
			w, errWorker := worker.DetailById(ctx, workId)
			if errWorker != nil || w == nil || w.Id == 0 {
				log.Warn("worker config is empty", zap.Any("worker info", w), zap.Error(errWorker))
				continue
			}

			common.JobList = w.JobList
		}
	}()
}

func Register(ctx context.Context) {
	errChan := make(chan error)
	go register(ctx, errChan)

	for {
		select {
		case workErr := <-errChan:
			log.Warn("worker routine stopped", zap.Error(workErr))
			go register(ctx, errChan)
		}
	}

}

func register(ctx context.Context, errChan chan error) {
	if workId == 0 {
		log.Panic("work id is empty")
	}

	//1.创建租约，设置过期时间
	leaseGrantResp, err := etcd.GetDefaultEtcdService().LeaseGrant(5)
	if err != nil {
		errChan <- err
		return
	}

	//2.创建上下文取消
	ctx, ctxFunc := context.WithCancel(context.Background())

	//3.保持租约通讯
	leaseKeepChan, err := etcd.GetDefaultEtcdService().LeaseKeepLive(ctx, leaseGrantResp.ID)
	if err != nil {
		_ = etcd.GetDefaultEtcdService().LeaseRevoke(leaseGrantResp.ID) //立刻释放租约
		ctxFunc()                                                       //停止自动续租
		errChan <- err
		return
	}

	errPut := etcd.GetDefaultEtcdService().PutWithTTLId("crontab", etcd.PrefixEnumWorker, fmt.Sprintf("%v", workId), "", leaseGrantResp.ID)
	if errPut != nil {
		errChan <- errPut
		return
	}

	//监听续租结果
	go func() {
		for {
			select {
			case keepResp := <-leaseKeepChan:
				if keepResp == nil {
					errChan <- errors.New("keep live fail")
					return
				}
			}
		}
	}()
}

func Jobs(ctx context.Context) {
	items, err := etcd.GetDefaultEtcdService().GetList("crontab", etcd.PrefixEnumJob, "list")
	if err != nil {
		log.Error("load grpc host list fail", zap.Error(err))
		return
	}
	for _, item := range items {
		var info *common2.JobInfo
		errJson := json.Unmarshal([]byte(item.Value[:]), &info)
		if errJson != nil {
			continue
		}
		job.PushJobChan(info)
	}

	go func() {
		for {
			log.Info("start watcher job")
			errWatch := jobWatcher()
			if errWatch != nil {
				log.Warn("job watcher stop ", zap.Error(errWatch))
			}
		}
	}()
}

func jobWatcher() (err error) {
	err = nil
	for {
		watchChan := etcd.GetDefaultEtcdService().WatchPath("crontab", etcd.PrefixEnumJob, "list")
		for watch := range watchChan {
			if watch.Canceled {
				return errors.New("watch closing")
			}
			for _, ev := range watch.Events {
				if ev.IsModify() || ev.IsCreate() {
					var info *common2.JobInfo
					errJson := json.Unmarshal(ev.Kv.Value[:], &info)
					if errJson != nil {
						continue
					}
					job.PushJobChan(info)
				}
			}
		}
	}
}
