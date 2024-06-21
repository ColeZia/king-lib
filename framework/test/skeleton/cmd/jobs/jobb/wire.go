//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"gl.king.im/king-lib/framework/test/skeleton/internal/conf"
	"gl.king.im/king-lib/framework/test/skeleton/internal/data"
	"gl.king.im/king-lib/framework/test/skeleton/internal/worker"

	//"gl.king.im/king-lib/protobuf/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/robfig/cron/v3"
)

// initApp init kratos application.
func initApp(*conf.Server, *conf.Data, *conf.Biz, log.Logger) (cron.Job, func(), error) {
	panic(wire.Build(data.ProviderSet, worker.ProviderJobBSet))
}
