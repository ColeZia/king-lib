package distr

import (
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/robfig/cron/v3"
	"gl.king.im/king-lib/framework/internal/data"
)

type DistrUpdateJobCnf struct {
	Close     bool   `protobuf:"varint,1,opt,name=close,proto3" json:"close,omitempty"`
	Spec      string `protobuf:"bytes,2,opt,name=spec,proto3" json:"spec,omitempty"`
	ReloadKey string `protobuf:"bytes,4,opt,name=reload_key,json=reloadKey,proto3" json:"reload_key,omitempty"`
}

const Key = ""

func WithRedisCli() {

}
func NewDistrUpdateJob(keyFunPairs []KeyFunPair, redisCli *redis.Client) (job cron.Job, err error) {

	if redisCli == nil {
		redisCli, err = data.GetRedisClient(nil)
		if err != nil {
			return
		}
	}

	job = &DistrUpdateJob{
		keyFunPairs: keyFunPairs,
		redisCli:    redisCli,
	}

	return
}

type KeyFunPair struct {
	Key           string
	Fun           func()
	lastUpdatedAt *time.Time
}
type DistrUpdateJob struct {
	redisCli      *redis.Client
	logger        log.Logger
	lastUpdatedAt *time.Time
	keyFunPairs   []KeyFunPair
}

// 注意这个如果不以指针作为接受者的话，lastUpdatedAt因为是局部变量就会一直nil
func (j *DistrUpdateJob) Run() {

	logFlag := "DistrUpdateJob.Run():"
	//log.Println(logFlag, "...doing...")

	for idx, pair := range j.keyFunPairs {
		pairLogFlag := logFlag + pair.Key + ":"
		retCmd := j.redisCli.Get(pair.Key)
		val, err := retCmd.Uint64()
		if err != nil {
			log.Println(pairLogFlag, "retCmd.Uint64() err:", err)
		} else {
			nowTime := time.Now()
			updatedAt := time.Unix(int64(val), 0)

			log.Println(pairLogFlag, "updatedAt:", updatedAt)
			if pair.lastUpdatedAt == nil || pair.lastUpdatedAt.Before(updatedAt) {
				//执行任务
				pair.Fun()
				j.keyFunPairs[idx].lastUpdatedAt = &nowTime
				log.Println(pairLogFlag, "pair.fun():", j.keyFunPairs[idx].lastUpdatedAt.String())
			}
		}
	}

	log.Println(logFlag, "...done")
}
