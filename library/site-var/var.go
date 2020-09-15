package site_var

import (
	"encoding/json"
	"github.com/garfieldlw/crontab-system/library/config"
	"github.com/garfieldlw/crontab-system/library/etcd"
	"sync"
)

var db = make(map[string]*DBConfig)
var lockDB = &sync.Mutex{}

type DBConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

func GetDBDefaultConfig() (*DBConfig, error) {
	env := config.GetENV()
	return GetDBConfig(env)
}

func GetDBConfig(env string) (*DBConfig, error) {
	if _, ok := db[env]; ok {
		return db[env], nil
	}

	lockDB.Lock()
	defer lockDB.Unlock()
	if _, ok := db[env]; !ok {
		val, err := etcd.GetDefaultEtcdService().GetBytes("common", etcd.PrefixEnumConfig, "db/"+env+"/common")
		if err != nil {
			return nil, err
		}

		var data *DBConfig
		errJson := json.Unmarshal(val, &data)
		if errJson != nil {
			return nil, errJson
		}

		db[env] = data
	}

	return db[env], nil
}
