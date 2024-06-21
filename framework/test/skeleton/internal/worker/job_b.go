package worker

import (
	"context"
	"log"
	"time"

	"github.com/google/wire"
	"github.com/robfig/cron/v3"
	"gl.king.im/king-lib/framework/test/skeleton/internal/data/ent/demo"
)

var ProviderJobBSet = wire.NewSet(NewJobB)

type JobB struct {
	dbObj  interface{}
	entCli *demo.Client
}

func (j JobB) Run() {
	log.Println("JobB.Run()...doing...")
	list, err := j.entCli.SomeTable.Query().All(context.Background())
	if err != nil {
		log.Println(err)
		return
	}

	for _, v := range list {
		log.Println(v)
	}
	time.Sleep(10 * time.Second)
	log.Println("JobB.Run()...done..")
}

func NewJobB(entCli *demo.Client) cron.Job {
	return JobB{entCli: entCli}
}
