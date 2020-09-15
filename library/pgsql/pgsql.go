package pgsql

import (
	"fmt"
	"github.com/garfieldlw/crontab-system/library/log"
	"github.com/garfieldlw/crontab-system/library/site-var"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"regexp"
	"sync"
)

type PgDbInfo struct {
	DbConfig *site_var.DBConfig
	Conn     *gorm.DB
}

var pgInstance *PgDbInfo
var lock = &sync.Mutex{}

func GetDb() *gorm.DB {
	pgInstance := LoadPgDb()
	if pgInstance == nil {
		return nil
	}
	return pgInstance.CheckAndReturnConn()
}

func LoadPgDb() *PgDbInfo {
	if pgInstance != nil {
		return pgInstance
	}

	lock.Lock()
	defer lock.Unlock()
	if pgInstance != nil {
		return pgInstance
	}

	dbConf, _ := site_var.GetDBDefaultConfig()

	PgDbInfo := new(PgDbInfo)
	PgDbInfo.DbConfig = dbConf

	PgDbInfo.InitConnect()

	pgInstance = PgDbInfo

	return pgInstance
}

func (PgDbInfo *PgDbInfo) InitConnect() {
	dbConf := PgDbInfo.DbConfig
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbConf.User, dbConf.Password, dbConf.Host, dbConf.Port, dbConf.Database)
	if db, e := gorm.Open("postgres", connStr); e != nil {
		log.Warn("load postage fail", zap.Error(e), zap.String("conn", connStr))
		PgDbInfo.Conn = nil
	} else {
		log.Warn("load postage success", zap.String("conn", connStr))
		db.DB().SetMaxIdleConns(5)
		db.DB().SetMaxOpenConns(20)
		PgDbInfo.Conn = db
	}
}

func (PgDbInfo *PgDbInfo) CheckAndReturnConn() *gorm.DB {
	if PgDbInfo.Conn == nil {
		lock.Lock()
		defer lock.Unlock()
		if PgDbInfo.Conn == nil {
			PgDbInfo.InitConnect()
		}
	}

	if err := PgDbInfo.Conn.DB().Ping(); err != nil {
		log.Warn("load postage fail", zap.Error(err))
		PgDbInfo.Clean()
		return nil
	}

	return PgDbInfo.Conn
}

func (PgDbInfo *PgDbInfo) Clean() {
	if PgDbInfo.Conn != nil {
		log.Warn("close postage conn")
		errClean := PgDbInfo.Conn.Close()
		if errClean != nil {
			log.Warn("close postage conn fail", zap.Error(errClean))
		}
	}

	PgDbInfo.Conn = nil
}

func EscapeNotFound(errGet error) error {
	if errGet != nil {
		reg := regexp.MustCompile(`not found`)
		finder := reg.FindAllString(errGet.Error(), -1)
		if len(finder) > 0 {
			return nil
		}
	}

	return errGet
}
