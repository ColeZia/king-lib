package main

import (
	"os"

	"gl.king.im/king-lib/framework/service"
	"gl.king.im/king-lib/framework/test/skeleton3/cmd/migrate"
	"gl.king.im/king-lib/framework/test/skeleton3/internal/conf"

	//_ "git.e.coding.king.cloud/dev/efficiency/open-lib/ci"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"gl.king.im/king-lib/framework/service/boot"
	v2 "gl.king.im/king-lib/framework/service/boot/v2"
	"gl.king.im/king-lib/framework/service/cmd"
)

var (
	// Name is the name of the compiled software.
	Name = "BossFrameworkTest3"
	// Version is the version of the compiled software.
	Version = "v0.0.1"
	// flagconf is the config flag.
	//flagconf string

	id, _ = os.Hostname()

	bc = &conf.Bootstrap{}
)

func newApp(logger log.Logger, sc *boot.NewServerCollection) *kratos.App {
	return v2.NewApp(logger, *sc)
}

func initAppWrapper(ll log.Logger) (*kratos.App, func(), error) {
	return initApp(bc.Server, bc.Data, bc.Biz, ll)
}

func initAppWrapperV2(bootData *service.ServiceBootData) (*kratos.App, func(), error) {
	return initApp(bc.Server, bc.Data, bc.Biz, bootData.Logger)
}

func main() {

	cmd.Execute(&cmd.ServiceInfo{
		Id:       id,
		Name:     Name,
		Version:  Version,
		BootConf: bc,
		//InitAppWrapper:   initAppWrapper,
		InitAppWrapperV2: initAppWrapperV2,
		MigrateRunFun:    migrate.MigrateRunFun,
		MigrateBc:        migrate.Bc,
	})

}
