package common

import "github.com/garfieldlw/crontab-system/library/etcd"

type ETCDKeysInputModel struct {
	Client string          `json:"client"`
	Prefix etcd.PrefixEnum `json:"prefix"`
	Key    string          `json:"key"`
}

type ETCDKeysOutputModel struct {
	Keys []*etcd.EtcdItem `json:"keys"`
}
