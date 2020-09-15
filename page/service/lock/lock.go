package lock

import (
	"github.com/garfieldlw/crontab-system/library/etcd"
	"context"
	"errors"
	"go.etcd.io/etcd/clientv3"
)

type Lock struct {
	lockKey    string
	cancelFunc context.CancelFunc
	leaseId    clientv3.LeaseID
	resFlag    bool
}

func CreateLock(client string, prefix etcd.PrefixEnum, key string) *Lock {
	return &Lock{
		lockKey: etcd.GetDefaultEtcdService().GetKey(client, prefix, key),
	}
}

//抢占锁
func (jobLock *Lock) Lock() error {
	//1.创建租约，设置过期时间
	leaseGrantResp, err := etcd.GetDefaultEtcdService().LeaseGrant(5)
	if err != nil {
		return err
	}

	//2.创建上下文取消
	ctx, ctxFunc := context.WithCancel(context.Background())

	//3.保持租约通讯
	_, errKeep := etcd.GetDefaultEtcdService().LeaseKeepLive(ctx, leaseGrantResp.ID)
	if errKeep != nil {
		_ = etcd.GetDefaultEtcdService().LeaseRevoke(leaseGrantResp.ID) //立刻释放租约
		ctxFunc()                                                       //停止自动续租
		return errKeep
	}

	//4.抢锁-创建事务
	txn := clientv3.KV(etcd.GetDefaultEtcdService().GetClient()).Txn(context.Background())
	//抢锁
	txn.If(clientv3.Compare(clientv3.CreateRevision(jobLock.lockKey), "=", 0)).
		Then(clientv3.OpPut(jobLock.lockKey, "", clientv3.WithLease(leaseGrantResp.ID))).
		Else(clientv3.OpGet(jobLock.lockKey))
	//5.提交事务
	txnResp, err := txn.Commit()
	if err != nil {
		_ = etcd.GetDefaultEtcdService().LeaseRevoke(leaseGrantResp.ID) //立刻释放租约
		ctxFunc()                                          //停止自动续租
		return err
	}
	//6.判断结果
	if !txnResp.Succeeded {
		//没有抢到锁
		_ = etcd.GetDefaultEtcdService().LeaseRevoke(leaseGrantResp.ID) //立刻释放租约
		ctxFunc()                                          //停止自动续租
		return errors.New("get lock fail")
	}

	//抢锁成功
	jobLock.resFlag = true
	jobLock.cancelFunc = ctxFunc
	jobLock.leaseId = leaseGrantResp.ID
	return nil

}

//取消锁
func (jobLock *Lock) Unlock() {
	if jobLock.resFlag {
		_ = etcd.GetDefaultEtcdService().LeaseRevoke(jobLock.leaseId)
		jobLock.cancelFunc()
	}
}
