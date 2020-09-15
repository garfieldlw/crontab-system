package config

import (
	"encoding/json"
	"github.com/garfieldlw/crontab-system/library/log"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"sync"
)

var info *Info
var once sync.Once

type Info struct {
	ENV string `json:"env"`
}

func GetENV() string {
	env := os.Getenv("ENVIRON")
	if len(env) > 0 {
		return env
	}

	if info == nil {
		readInfo()
	}

	return info.ENV
}

func readInfo() {
	once.Do(func() {
		data, err := ReadJsonFile("config/config.json")
		if err != nil {
			log.Fatal("load config err", zap.Error(err))
			os.Exit(1)
		}

		errJson := json.Unmarshal(data, &info)
		if errJson != nil {
			log.Fatal("load config err", zap.Error(errJson))
			os.Exit(1)
		}
	})
}

func ReadJsonFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}
