package worker

import (
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/robfig/cron/v3"
	"gl.king.im/king-lib/framework/scheduler/distr"
)

func NewDistrUpdateJob() cron.Job {

	job, err := distr.NewDistrUpdateJob([]distr.KeyFunPair{{
		Key: "reload_policy",
		Fun: func() {
			log.Println("reload_policy...")
		},
	}, {
		Key: "reload_policy222",
		Fun: func() {
			log.Println("reload_policy222...")
		},
	},
	}, nil)

	if err != nil {
		panic(err)
	}

	return job
}

type DistrUpdateJob struct {
	reloadKey        string
	redisCli         *redis.Client
	logger           log.Logger
	lastUpdatedAt    *time.Time
	extLastUpdatedAt *time.Time
}
