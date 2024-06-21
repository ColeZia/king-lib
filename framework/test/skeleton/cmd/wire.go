//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"gl.king.im/king-lib/framework/test/skeleton/internal/biz"
	"gl.king.im/king-lib/framework/test/skeleton/internal/conf"
	"gl.king.im/king-lib/framework/test/skeleton/internal/data"
	"gl.king.im/king-lib/framework/test/skeleton/internal/server"
	"gl.king.im/king-lib/framework/test/skeleton/internal/service"
	"gl.king.im/king-lib/framework/test/skeleton/internal/worker"

	//"gl.king.im/king-lib/protobuf/conf"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// initApp init kratos application.
func initApp(*conf.Server, *conf.Data, *conf.Biz, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, worker.ProviderSet, newApp))
}
