package worker

import (
	"time"

	"github.com/google/wire"
	"gl.king.im/king-lib/framework/scheduler"
	"gl.king.im/king-lib/framework/test/skeleton2/internal/data/ent/demo"
)

var ProviderSet = wire.NewSet(NewScheduler)

func NewScheduler(entCli *demo.Client) (sche scheduler.DistributedScheduler, cleanup func(), err error) {
	sche, err = scheduler.NewDistributedScheduler([]*scheduler.Job{
		{
			Close:      true,
			Spec:       "@every 20s",
			CronJob:    JobA{entCli: entCli},
			RunOnStart: true,
			//计划任务的执行时间可能超过默认的互斥锁过期时长时，需要主动设置锁的过期时长，默认过期时长可查看scheduler.DefaultMutexDuration常量值
			LockExpiredDuration: 24 * time.Hour,
		},
		{
			Close:   true,
			Spec:    "@every 2s",
			CronJob: &JobB{},
		},
		{
			Close:   true,
			Spec:    "@every 2s",
			CronJob: NewDistrUpdateJob(),
		},
	})

	if err != nil {
		panic(err)
	}

	sche.Start()

	cleanup = func() {
		sche.Cleanup()
	}

	return
}
