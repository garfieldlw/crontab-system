package etcd

import (
	"context"
	"errors"
	"fmt"
	"github.com/garfieldlw/crontab-system/library/config"
	"github.com/garfieldlw/crontab-system/library/log"
	"go.etcd.io/etcd/clientv3"
	"go.uber.org/zap"
	"regexp"
	"sync"
	"time"
)

const TIMEOUT = 5 * time.Second

var lock = &sync.Mutex{}
var etcdInstance *Etcd
var env string

type EtcdConf struct {
	Urls     []string
	UserName string
	Password string
}

var productionEtcdConf = &EtcdConf{
	Urls: []string{

	},
}

var developmentEtcdConf = &EtcdConf{
	Urls: []string{
		"127.0.0.1:4001",
	},
	UserName: "username",
	Password: "password",
}

type Etcd struct {
	cli *clientv3.Client
}

type EtcdItem struct {
	Path  string `json:"path"`
	Value string `json:"value"`
}

func GetDefaultEtcdService() *Etcd {
	if etcdInstance != nil {
		return etcdInstance
	}

	lock.Lock()
	defer lock.Unlock()

	if etcdInstance == nil {
		var etcdConfig *EtcdConf
		switch env {
		case "production":
			etcdConfig = productionEtcdConf
		default:
			etcdConfig = developmentEtcdConf
		}
		etcdInstance = initEtcd(etcdConfig)
	}

	return etcdInstance
}

func initEtcd(etcdConfig *EtcdConf) *Etcd {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:            etcdConfig.Urls,
		DialTimeout:          TIMEOUT,
		PermitWithoutStream:  true,
		DialKeepAliveTime:    TIMEOUT,
		DialKeepAliveTimeout: TIMEOUT,
	})

	if err != nil {
		log.Panic("connect etcd failed", zap.Error(err))
	}

	return &Etcd{cli: cli}
}

func (etcd *Etcd) CheckKey(client string, prefix PrefixEnum, keyPath string) (bool, error) {
	keyPath = etcd.GetKey(client, prefix, keyPath)
	kv := clientv3.KV(etcd.cli)
	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT)
	res, err := kv.Get(ctx, keyPath)
	if err != nil {
		return false, err
	}

	if res.Count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (etcd *Etcd) Delete(client string, prefix PrefixEnum, keyPath string) error {
	keyPath = etcd.GetKey(client, prefix, keyPath)
	kv := clientv3.KV(etcd.cli)
	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT)
	res, err := kv.Delete(ctx, keyPath)
	log.Info("delete key", zap.String("keyPath", keyPath), zap.Any("res", res), zap.Error(err))
	if err != nil {
		return err
	}
	return nil
}

func (etcd *Etcd) DeleteByKey(key string) error {
	kv := clientv3.KV(etcd.cli)
	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT)
	res, err := kv.Delete(ctx, key)
	log.Info("delete key", zap.String("keyPath", key), zap.Any("res", res), zap.Error(err))
	if err != nil {
		return err
	}
	return nil
}

func (etcd *Etcd) Put(client string, prefix PrefixEnum, keyPath, value string) error {
	keyPath = etcd.GetKey(client, prefix, keyPath)
	kv := clientv3.KV(etcd.cli)
	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT)
	res, err := kv.Put(ctx, keyPath, value)
	log.Info("put value", zap.String("keyPath", keyPath), zap.Any("res", res), zap.Error(err))
	if err != nil {
		return err
	}
	return nil
}

func (etcd *Etcd) PutWithTTL(client string, prefix PrefixEnum, keyPath, value string, ttl int64) error {
	keyPath = etcd.GetKey(client, prefix, keyPath)
	respTTL, errTTL := etcd.LeaseGrant(ttl)
	if errTTL != nil {
		return errTTL
	}

	kv := clientv3.KV(etcd.cli)
	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT)
	res, err := kv.Put(ctx, keyPath, value, clientv3.WithLease(respTTL.ID))
	log.Info("put value with ttl", zap.String("keyPath", keyPath), zap.Any("res", res), zap.Error(err))
	if err != nil {
		return err
	}
	return nil
}

func (etcd *Etcd) PutWithTTLId(client string, prefix PrefixEnum, keyPath, value string, ttlId clientv3.LeaseID) error {
	keyPath = etcd.GetKey(client, prefix, keyPath)
	kv := clientv3.KV(etcd.cli)
	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT)
	res, err := kv.Put(ctx, keyPath, value, clientv3.WithLease(ttlId))
	log.Info("put value with ttl", zap.String("keyPath", keyPath), zap.Any("res", res), zap.Error(err))
	if err != nil {
		return err
	}
	return nil
}

func (etcd *Etcd) Get(client string, prefix PrefixEnum, keyPath string) (string, error) {
	keyPath = etcd.GetKey(client, prefix, keyPath)
	kv := clientv3.KV(etcd.cli)
	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT)
	res, err := kv.Get(ctx, keyPath)
	if err != nil {
		return "", err
	}

	for _, val := range res.Kvs {
		if string(val.Key[:]) == keyPath {
			return string(val.Value[:]), nil
		}
	}

	return "", errors.New("no value in etcd")
}

func (etcd *Etcd) GetBytes(client string, prefix PrefixEnum, keyPath string) ([]byte, error) {
	keyPath = etcd.GetKey(client, prefix, keyPath)
	kv := clientv3.KV(etcd.cli)
	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT)
	res, err := kv.Get(ctx, keyPath)
	if err != nil {
		return nil, err
	}

	for _, val := range res.Kvs {
		if string(val.Key[:]) == keyPath {
			return val.Value, nil
		}
	}

	return nil, errors.New("no value in etcd")
}

func (etcd *Etcd) GetList(client string, prefix PrefixEnum, keyPath string) ([]*EtcdItem, error) {
	keyPath = etcd.GetKey(client, prefix, keyPath)
	kv := clientv3.KV(etcd.cli)
	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT)
	res, err := kv.Get(ctx, keyPath, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	items := make([]*EtcdItem, 0)

	for _, val := range res.Kvs {
		item := new(EtcdItem)
		item.Path = string(val.Key[:])
		item.Value = string(val.Value[:])

		items = append(items, item)
	}

	return items, nil
}

func (etcd *Etcd) LeaseGrant(ttl int64) (*clientv3.LeaseGrantResponse, error) {
	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT)
	lease := clientv3.Lease(etcd.cli)
	res, err := lease.Grant(ctx, ttl)
	if err != nil {
		log.Info("lease grant", zap.Any("res", res), zap.Error(err))
		return nil, err
	}
	return res, nil
}

func (etcd *Etcd) LeaseKeepLive(ctx context.Context, id clientv3.LeaseID) (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	lease := clientv3.Lease(etcd.cli)
	res, err := lease.KeepAlive(ctx, id)
	if err != nil {
		log.Info("lease grant", zap.Any("res", res), zap.Error(err))
		return nil, err
	}
	return res, nil
}

func (etcd *Etcd) LeaseKeepLiveOnce(ctx context.Context, id clientv3.LeaseID) (*clientv3.LeaseKeepAliveResponse, error) {
	lease := clientv3.Lease(etcd.cli)
	res, err := lease.KeepAliveOnce(ctx, id)
	if err != nil {
		log.Info("lease grant", zap.Any("res", res), zap.Error(err))
		return nil, err
	}
	return res, nil
}

func (etcd *Etcd) LeaseRevoke(id clientv3.LeaseID) error {
	ctx, _ := context.WithTimeout(context.Background(), TIMEOUT)
	lease := clientv3.Lease(etcd.cli)
	res, err := lease.Revoke(ctx, id)
	if err != nil {
		log.Info("lease revoke", zap.Any("res", res), zap.Error(err))
		return err
	}
	return nil
}

func (etcd *Etcd) WatchKey(client string, prefix PrefixEnum, key string) clientv3.WatchChan {
	key = etcd.GetKey(client, prefix, key)
	return etcd.cli.Watch(context.Background(), key)
}

func (etcd *Etcd) WatchPath(client string, prefix PrefixEnum, keyPath string) clientv3.WatchChan {
	keyPath = etcd.GetKey(client, prefix, keyPath)
	return etcd.cli.Watch(context.Background(), keyPath, clientv3.WithPrefix())
}

func (etcd *Etcd) GetClient() *clientv3.Client {
	return etcd.cli
}

func init() {
	env = config.GetENV()
}

func (etcd *Etcd) GetKey(client string, prefix PrefixEnum, key string) string {
	if len(client) == 0 {
		client = "common"
	}
	return fmt.Sprintf("/%s/%s/%s", client, string(prefix), key)
}

func (etcd *Etcd) EscapeUnauthenticated(errGet error) error {
	if errGet != nil {
		reg := regexp.MustCompile(`invalid auth token`)
		finder := reg.FindAllString(errGet.Error(), -1)
		if len(finder) > 0 {
			return nil
		}
	}

	return errGet
}

type PrefixEnum string

const (
	PrefixEnumConfig PrefixEnum = "config"
	PrefixEnumWorker PrefixEnum = "worker"
	PrefixEnumMaster PrefixEnum = "master"
	PrefixEnumJob    PrefixEnum = "job"
	PrefixEnumFlow   PrefixEnum = "flow"
)
