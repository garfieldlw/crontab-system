package unique_id

import (
	"github.com/garfieldlw/crontab-system/library/log"
	"github.com/garfieldlw/crontab-system/library/unique"
	"errors"
	"go.uber.org/zap"
)

func GetUniqueId(bizType unique.BizType) (int64, error) {
	ids, idErr := unique.GenerateUniqueId(bizType, 1)
	if idErr != nil {
		log.Warn("get unique id error", zap.Error(idErr))
		return 0, idErr
	}

	if len(ids) == 0 {
		return 0, errors.New("get id fail")
	}

	return ids[0], nil
}

func GetUniqueIds(bizType unique.BizType, num int32) ([]int64, error) {
	ids, idErr := unique.GenerateUniqueId(bizType, num)
	if idErr != nil {
		log.Warn("get unique id error", zap.Error(idErr))
		return nil, idErr
	}

	return ids, nil
}
