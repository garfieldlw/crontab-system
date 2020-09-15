package unique

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

const (
	twepoch        = int64(1546272000000)             // 开始时间截 (2019-01-01)
	workerBits     = uint(10)                         // 机器id所占的位数
	typeBits       = uint(7)                          // type所占位数
	sequenceBits   = uint(7)                          // 序列所占的位数
	sequenceMask   = int64(-1 ^ (-1 << sequenceBits)) //
	sequenceShift  = uint(0)                          // 序列左移位数
	workerShift    = sequenceBits                     // 机器id左移位数
	timestampShift = sequenceBits + workerBits        // 时间戳左移位数
	typeShift      = uint(63) - typeBits
)

var (
	service *Snowflake
	lock    = &sync.Mutex{}
)

type Snowflake struct {
	sync.Mutex
	timestamp int64
	worker    int64
	sequence  int64
}

type IdInfo struct {
	BizType   BizType `json:"biz_type"`
	Timestamp int64   `json:"timestamp"`
}

func init() {
	if service != nil {
		return
	}

	lock.Lock()
	defer lock.Unlock()
	if service == nil {
		s, err := NewSnowflake()
		if err != nil {
			log.Fatal(err)
		}

		service = s
	}
}

func GenerateUniqueId(bizType BizType, num int32) ([]int64, error) {
	result := make([]int64, 0)
	var a int32 = 0
	for a = 0; a < num; a++ {
		result = append(result, service.Generate(bizType))
	}

	return result, nil
}

func GetUniqueIdInfo(id int64) (*IdInfo, error) {
	info := new(IdInfo)

	info.BizType = getBizType(id)
	info.Timestamp = getTimestamp(id)

	return info, nil
}

func NewSnowflake() (*Snowflake, error) {
	workId, errWorkId := getWorkId()
	if errWorkId != nil {
		return nil, errWorkId
	}

	return &Snowflake{
		timestamp: 0,
		worker:    workId,
		sequence:  0,
	}, nil
}

func (s *Snowflake) Generate(bizType BizType) int64 {
	s.Lock()
	defer s.Unlock()

	now := time.Now().UnixNano() / 1000000

	if s.timestamp == now {
		s.sequence = (s.sequence + 1) & sequenceMask

		if s.sequence == 0 {
			for now <= s.timestamp {
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		s.sequence = 0
	}

	s.timestamp = now

	return int64(((now-twepoch)&0x7FFFFFFFFF)<<timestampShift | (s.worker << workerShift) | (s.sequence << sequenceShift) | int64(bizType)<<typeShift)
}

func getWorkId() (int64, error) {
	workerId := hostnameToInt() & 0x3FF
	if workerId < 0 || workerId > 0x3FF {
		return 0, errors.New("get work id error")
	}

	return int64(workerId), nil
}

func getBizType(id int64) BizType {
	bizType := (id >> typeShift) & 0x7F
	return BizType(int32(bizType))
}

func getTimestamp(id int64) int64 {
	ts := (id >> timestampShift) & 0x7FFFFFFFFF
	return ts + twepoch
}

func hostnameToInt() uint {
	var hostname, err = os.Hostname()
	var hash = uint(0)
	var seed = uint(131)
	if err != nil {
		return 0
	}

	fmt.Println(hostname)

	for _, k := range hostname {
		hash = hash*seed + uint(k)
	}
	fmt.Println(hash)

	return hash
}
