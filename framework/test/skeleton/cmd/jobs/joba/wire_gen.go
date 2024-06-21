// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2/log"
	"gl.king.im/king-lib/framework/scheduler"
	"gl.king.im/king-lib/framework/test/skeleton/internal/conf"
	"gl.king.im/king-lib/framework/test/skeleton/internal/data"
	"gl.king.im/king-lib/framework/test/skeleton/internal/worker"
)

// Injectors from wire.go:

// initApp init kratos application.
func initApp(server *conf.Server, confData *conf.Data, biz *conf.Biz, logger log.Logger) (scheduler.DistributedScheduler, func(), error) {
	client := data.NewBossEntCli(confData)
	distributedScheduler, cleanup, err := worker.NewScheduler(client)
	if err != nil {
		return nil, nil, err
	}
	return distributedScheduler, func() {
		cleanup()
	}, nil
}
