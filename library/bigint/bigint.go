package bigint

import (
	"encoding/json"
	"github.com/garfieldlw/crontab-system/library/log"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type BigInt struct {
	str string
}

func (v *BigInt) ToString() string {
	return v.str
}

func (v *BigInt) ToInt64() int64 {
	return cast.ToInt64(v.str)
}

func (v *BigInt) MarshalJSON() ([]byte, error) {
	value, errJson := json.Marshal(v.str)
	if errJson != nil {
		log.Error("big int marshal error", zap.Error(errJson))
	}
	return value, errJson
}

func (v *BigInt) UnmarshalJSON(b []byte) error {
	var val string
	errJson := json.Unmarshal(b, &val)
	if errJson != nil {
		log.Error("big int unmarshal error", zap.String("origin data", string(b[:])), zap.String("error info", errJson.Error()))
		v.str = string(b[:])
	} else {
		v.str = val
	}

	return nil
}
