package etcd

import (
	"github.com/garfieldlw/crontab-system/library/etcd"
	"github.com/garfieldlw/crontab-system/page/service/common"
	"context"
)

func GetKeys(ctx context.Context, client string, prefix etcd.PrefixEnum, key string) (*common.ETCDKeysOutputModel, error) {
	items, err := etcd.GetDefaultEtcdService().GetList(client, prefix, key)
	if err != nil {
		return nil, err
	}

	return &common.ETCDKeysOutputModel{
		Keys: items,
	}, nil
}
