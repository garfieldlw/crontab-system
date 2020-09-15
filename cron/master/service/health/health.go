package health

import (
	"context"
	"fmt"
	"github.com/garfieldlw/crontab-system/library/etcd"
	"github.com/garfieldlw/crontab-system/library/log"
	"github.com/garfieldlw/crontab-system/library/message/email"
	"github.com/garfieldlw/crontab-system/page/service/worker"
	"go.uber.org/zap"
	"sync"
	"time"
)

var (
	checkWorkerList     []int64
	checkWorkerMap      = &sync.Map{}
	checkWorkerAlertMap = &sync.Map{}
)

func CheckHealth(ctx context.Context) {
	for _, id := range checkWorkerList {
		checkWorkerMap.Store(id, int64(0))
		checkWorkerAlertMap.Store(id, int32(0))
	}
	go checkHealth(ctx)
}

func checkHealth(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 60 * 5)
	//ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	for {
		<-ticker.C
		for _, id := range checkWorkerList {
			has, _ := etcd.GetDefaultEtcdService().CheckKey("crontab", etcd.PrefixEnumWorker, fmt.Sprintf("%v", id))
			if has {
				checkWorkerMap.Store(id, time.Now().Unix())
				checkWorkerAlertMap.Store(id, 0)
				continue
			}

			need, count := needAlert(id)
			if !need {
				continue
			}

			checkWorkerAlertMap.Store(id, count)

			w, errWorker := worker.DetailById(ctx, id)
			if errWorker != nil || w == nil || w.Id == 0 {
				continue
			}

			log.Warn("worker offline", zap.Int32("count", count), zap.Any("worker info", w))

			_ = email.SendMailTLS("username@email.com", "工作流处理失败[服务器不在线]", fmt.Sprintf("服务器不在线，服务器id：%v，服务器名称：%v", w.Id, w.Name), "")
		}
	}
}

func needAlert(id int64) (bool, int32) {
	if val, ok := checkWorkerMap.Load(id); ok {
		if val == nil {
			return true, 1
		}

		if valx, okValue := val.(int64); okValue {
			x := time.Now().Unix() - valx
			if x < 60*10 {
				return false, 0
			}

			if count, okCount := checkWorkerAlertMap.Load(id); okCount {
				if count == nil {
					return true, 1
				}
				if countVal, okCountVal := count.(int32); okCountVal {
					if countVal > 3 {
						return false, countVal
					} else {
						return true, countVal + 1
					}
				} else {
					return true, 1
				}
			} else {
				return true, 1
			}
		} else {
			return true, 1
		}
	} else {
		return true, 1
	}
}
