package scheduler

import (
	"errors"
	"log"
	"reflect"
	"time"

	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis"
	"github.com/robfig/cron/v3"
	"gl.king.im/king-lib/framework/config"
	"gl.king.im/king-lib/framework/internal/data"
	"gl.king.im/king-lib/framework/service"
)

type DistributedScheduler interface {
	Start()
	Stop()
	Cleanup()
	//GetCron() *cron.Cron
	//生成一个互斥锁
	//GenrateJobMutex(jobName string, lockExpiredDuration time.Duration) (mutex *redsync.Mutex, err error)
	//FetchJob(jobName string) (info Job, err error)
	//Lock方法会生成一个redsync.Mutex实例，同时执行该实例的Lock方法，再返回该实例
	//Lock(jobName string) (mutex *redsync.Mutex, err error)
	//Unlock() error
}

type Job struct {
	Close      bool
	RunOnStart bool
	//请注意JobName此名称将作为互斥锁的key
	JobName string
	CronJob cron.Job
	Spec    string
	//互斥锁的过期时长，默认时长查看DefaultMutexDuration常量值
	LockExpiredDuration time.Duration
	mutex               *redsync.Mutex
	finalJobName        string
	//Locked              bool
	//Locker              string
}

type SchedulerImpl struct {
	opts    *SchedulerOptions
	cronIns *cron.Cron
	jobs    []*Job
	//用于不同服务间的任务命名隔离，默认为服务名
	namespace string
}

type SchedulerOption func(*SchedulerOptions)

const defaultMutexPrefix = "SchedulerMutex:"
const DefaultMutexDuration = time.Minute * 30

type SchedulerOptions struct {
	redisCli    *redis.Client
	jobs        []*Job
	rdsCnf      *config.Data_Redis
	sync        *redsync.Redsync
	mutexPrefix string
	logger      klog.Logger
}

func WithRedisCli(opt *redis.Client) SchedulerOption {
	return func(o *SchedulerOptions) {
		o.redisCli = opt
	}
}

func WithRedisConf(opt *config.Data_Redis) SchedulerOption {
	return func(o *SchedulerOptions) {
		o.rdsCnf = opt
	}
}

func WithLogger(opt klog.Logger) SchedulerOption {
	return func(o *SchedulerOptions) {
		o.logger = opt
	}
}

//func WithJobs(opt []*Job) SchedulerOption {
//	return func(o *SchedulerOptions) {
//		o.jobs = opt
//	}
//}

type cronLogger struct {
}

func (cronLogger) Info(msg string, keysAndValues ...interface{}) {
	log.Println("cronLogger::", msg, keysAndValues)
}

func (cronLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	log.Println("cronLogger::", err, msg, keysAndValues)
}

func NewDistributedScheduler(jobs []*Job, opts ...SchedulerOption) (sche DistributedScheduler, err error) {

	finalOpts := &SchedulerOptions{
		mutexPrefix: defaultMutexPrefix,
	}

	for _, v := range opts {
		v(finalOpts)
	}

	if finalOpts.redisCli == nil {
		var cli *redis.Client
		cli, err = data.GetRedisClient(finalOpts.rdsCnf)
		if err != nil {
			return
		}
		finalOpts.redisCli = cli
	}

	pool := goredis.NewPool(finalOpts.redisCli)

	rs := redsync.New(pool)

	finalOpts.sync = rs

	scheImpl := &SchedulerImpl{
		opts: finalOpts,
	}

	if len(jobs) > 0 {
		scheImpl.cronIns = cron.New()

		for k, v := range jobs {

			jobName := v.JobName
			if jobName == "" {
				cronJobType := reflect.TypeOf(v.CronJob)
				switch cronJobType.Kind() {
				case reflect.Ptr:
					jobName = cronJobType.Elem().Name()
				default:
					jobName = cronJobType.Name()
				}

			}
			jobs[k].JobName = jobName

			if v.Close {
				log.Println("notice:" + jobName + "任务未开启")
				continue
			}

			if v.CronJob == nil {
				err = errors.New("计划任务不能为空!")
				return
			}

			mutexDur := DefaultMutexDuration
			if v.LockExpiredDuration > 0 {
				mutexDur = v.LockExpiredDuration
			}

			mutexKey := scheImpl.opts.mutexPrefix + service.AppInfoIns.Name + ":" + jobName
			jobs[k].mutex = finalOpts.sync.NewMutex(mutexKey, redsync.WithExpiry(mutexDur))

			cronChainJob := cron.NewChain(
				cron.SkipIfStillRunning(cronLogger{}),
				//注意这里不能传入v
				SkipIfDistributedMutexLocked(jobs[k]),
			).Then(v.CronJob)

			var entryID cron.EntryID
			entryID, err = scheImpl.cronIns.AddJob(v.Spec, cronChainJob)
			if err != nil {
				log.Println(jobName + "任务添加失败:" + err.Error())
				return
			}

			_ = entryID

			log.Println(jobName + "任务已添加，spec: " + v.Spec + "; mutexKey: " + mutexKey)

			if v.RunOnStart {
				log.Println(jobName + "任务设置为启动时执行一次")
				jobWrapper := distrRunner(v.CronJob, v)
				go jobWrapper()
			}

		}
		scheImpl.jobs = jobs
	} else {
		log.Println("计划任务未开启")
	}

	sche = scheImpl

	return
}

// 分布式锁实现
func SkipIfDistributedMutexLocked(job *Job) cron.JobWrapper {
	return func(j cron.Job) cron.Job {
		return cron.FuncJob(distrRunner(j, job))
	}
}

func distrRunner(j cron.Job, job *Job) func() {
	return func() {
		if err := job.mutex.Lock(); err != nil {
			log.Println("job:"+job.JobName+": job.mutex.Lock() err:", err)
			return
		}

		//log.Println("locked working start...")
		j.Run()
		//log.Println("locked working over...")

		if ok, err := job.mutex.Unlock(); !ok || err != nil {
			log.Println("job:"+job.JobName+": job.mutex.Unlock() err:", err)
			return
		}
	}
}

func (s *SchedulerImpl) Start() {
	s.cronIns.Start()
}

func (s *SchedulerImpl) Stop() {
	s.cronIns.Stop()
}

func (s *SchedulerImpl) Cleanup() {
	s.cronIns.Stop()
	for _, v := range s.jobs {
		if v.mutex == nil {
			continue
		}

		if ok, err := v.mutex.Unlock(); !ok || err != nil {
			log.Println("任务:"+v.JobName+" 锁释放失败或已释放过:", err)
			return
		} else {
			log.Println("任务:" + v.JobName + " 锁释放成功")
		}
	}

	log.Println("Scheduler Cleanup done")
}

func (s *SchedulerImpl) getCron() *cron.Cron {
	return s.cronIns
}

func (s *SchedulerImpl) GenrateJobMutex(jobName string, expired time.Duration) (mutex *redsync.Mutex, err error) {

	mutexname := jobName
	mutex = s.opts.sync.NewMutex(mutexname, redsync.WithExpiry(expired))
	return
}
