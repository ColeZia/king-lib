package worker

import (
	"context"
	"log"
	"time"

	"gl.king.im/king-lib/framework/test/skeleton2/internal/data/ent/demo"
)

type JobA struct {
	entCli *demo.Client
}

func (j JobA) Run() {
	log.Println("JobA.Run()...doing...")

	list, err := j.entCli.SomeTable.Query().All(context.Background())
	if err != nil {
		log.Println(err)
		return
	}

	for _, v := range list {
		log.Println(v)
	}

	_ = list
	time.Sleep(10 * time.Second)

	log.Println("JobA.Run()...done..")
}
