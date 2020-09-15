package init_cron

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/garfieldlw/crontab-system/cron/master/service/schedule"
	"github.com/garfieldlw/crontab-system/library/etcd"
	"github.com/garfieldlw/crontab-system/library/log"
	common2 "github.com/garfieldlw/crontab-system/page/service/common"
	"go.uber.org/zap"
)

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

	errPut := etcd.GetDefaultEtcdService().PutWithTTLId("crontab", etcd.PrefixEnumMaster, "master", "", leaseGrantResp.ID)
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

func Flows(ctx context.Context) {
	items, err := etcd.GetDefaultEtcdService().GetList("crontab", etcd.PrefixEnumFlow, "list")
	if err != nil {
		log.Error("load grpc host list fail", zap.Error(err))
		return
	}
	for _, item := range items {
		var info *common2.FlowInfo
		errJson := json.Unmarshal([]byte(item.Value[:]), &info)
		if errJson != nil {
			continue
		}
		schedule.PushToFlowChan(info)
	}

	go func() {
		for {
			log.Info("start watcher flow")
			errWatch := flowWatcher()
			if errWatch != nil {
				log.Warn("job watcher stop ", zap.Error(errWatch))
			}
		}
	}()
}

func flowWatcher() (err error) {
	err = nil
	for {
		watchChan := etcd.GetDefaultEtcdService().WatchPath("crontab", etcd.PrefixEnumFlow, "list")
		for watch := range watchChan {
			if watch.Canceled {
				return errors.New("watch closing")
			}
			for _, ev := range watch.Events {
				if ev.IsModify() || ev.IsCreate() {
					var info *common2.FlowInfo
					errJson := json.Unmarshal(ev.Kv.Value[:], &info)
					if errJson != nil {
						continue
					}
					schedule.PushToFlowChan(info)
				}
			}
		}
	}
}
