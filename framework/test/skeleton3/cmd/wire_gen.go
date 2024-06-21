// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"gl.king.im/king-lib/framework/test/skeleton3/internal/biz"
	"gl.king.im/king-lib/framework/test/skeleton3/internal/conf"
	"gl.king.im/king-lib/framework/test/skeleton3/internal/data"
	"gl.king.im/king-lib/framework/test/skeleton3/internal/server"
	v1_2 "gl.king.im/king-lib/framework/test/skeleton3/internal/service/admin/v1"
	"gl.king.im/king-lib/framework/test/skeleton3/internal/service/service/v1"
	"gl.king.im/king-lib/framework/test/skeleton3/internal/worker"
)

// Injectors from wire.go:

// initApp init kratos application.
func initApp(confServer *conf.Server, confData *conf.Data, confBiz *conf.Biz, logger log.Logger) (*kratos.App, func(), error) {
	bossDb := data.NewBossDb(confData)
	client := data.NewBossEntCli(confData)
	dataData, cleanup, err := data.NewData(confData, logger, bossDb, client)
	if err != nil {
		return nil, nil, err
	}
	skeletonRepo := data.NewSkeletonRepo(dataData, logger)
	skeletonUsecase := biz.NewSkeletonUsecase(skeletonRepo, logger, confBiz)
	skeletonService := v1.NewSkeletonService(skeletonUsecase)
	skeletonGinService := v1.NewSkeletonGinService(skeletonUsecase)
	skeletonAdminService := v1_2.NewSkeletonAdminService(skeletonUsecase)
	distributedScheduler, cleanup2, err := worker.NewScheduler(client)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	newServerCollection := server.NewServer(skeletonService, skeletonGinService, skeletonAdminService, distributedScheduler, logger)
	app := newApp(logger, newServerCollection)
	return app, func() {
		cleanup2()
		cleanup()
	}, nil
}
