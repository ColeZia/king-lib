package main

import (
	"os"

	//"gl.king.im/king-lib/protobuf/conf"

	"github.com/go-kratos/kratos/v2/log"
	"gl.king.im/king-lib/framework/scheduler"
	"gl.king.im/king-lib/framework/service"
	"gl.king.im/king-lib/framework/test/skeleton/internal/conf"

	fcmd "gl.king.im/king-lib/framework/service/cmd/scheduler"
	"gl.king.im/king-lib/framework/test/skeleton/cmd/migrate"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()

	bc = &conf.Bootstrap{}
)

func init() {
	Name = "BossOrder"
	Version = "v0.0.21"
}

func initAppWrapper(lo log.Logger) (scheduler.DistributedScheduler, func(), error) {
	return initApp(bc.Server, bc.Data, bc.Biz, lo)
}

func initAppWrapperV2(bootData *service.ServiceBootData) (scheduler.DistributedScheduler, func(), error) {
	return initApp(bc.Server, bc.Data, bc.Biz, bootData.Logger)
}

func main() {
	//flag.Parse()
	fcmd.Execute(&fcmd.ServiceInfo{
		Id:              id,
		Name:            Name,
		Version:         Version,
		BootConf:        bc,
		MigrateRunFun:   migrate.MigrateRunFun,
		MigrateBc:       migrate.Bc,
		InitScheWrapper: initAppWrapperV2,
	})
}
