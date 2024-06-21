/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package job

import (
	"gl.king.im/king-lib/framework/auth/user"
	"gl.king.im/king-lib/framework/internal/di"
	"gl.king.im/king-lib/framework/internal/stat"
	"gl.king.im/king-lib/framework/log"
	"gl.king.im/king-lib/framework/service"
	"gl.king.im/king-lib/framework/service/cmd/config"

	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/spf13/cobra"
	fwConf "gl.king.im/king-lib/framework/config"
)

// serviceCmd represents the service command
var jobCmd = &cobra.Command{
	Use:   "service",
	Short: "微服务相关命令",
	Long:  `微服务相关命令，包括启动微服务等功能`,
	Run: func(cmd *cobra.Command, args []string) {
		runJob()
	},
}

// serviceCmd represents the service command
var jobRunCmd = &cobra.Command{
	Use:   "run",
	Short: "微服务启动命令",
	Long:  `启动微服务`,
	Run: func(cmd *cobra.Command, args []string) {
		runJob()
	},
}

func init() {
	jobRunCmd.PersistentFlags().StringVarP(&confpath, "conf", "c", "", "配置文件")
	jobCmd.AddCommand(jobRunCmd)

	rootCmd.AddCommand(jobCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//serviceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func runJob() {
	config.ScanCnf(srvInfo.Name, confpath, srvInfo.BootConf)

	logger, err := log.NewLogger("default")
	if err != nil {
		panic(err)
	}

	logger = klog.With(logger,
		//"ts", klog.DefaultTimestamp,
		//"caller", klog.DefaultCaller,
		"service.id", srvInfo.Id,
		"service.name", srvInfo.Name,
		"service.version", srvInfo.Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
		"user.username", user.UsernameLogValuer(),
		"operation", log.OperationValuer(),
		"short_op", log.ShortOperationValuer(),
		//"http", log.HttpValuer(),
		"http.method", log.HttpMethodValuer(),
		"http.host", log.HttpHostValuer(),
		"http.url.path", log.UrlPathValuer(),
		"duration", stat.StatDurationStrValuer(),
	)

	di.SetLogger(logger)

	jobIns, cleanup, err := srvInfo.InitJobWrapper(&service.ServiceBootData{
		Name:     srvInfo.Name,
		Version:  srvInfo.Version,
		Id:       srvInfo.Id,
		Logger:   logger,
		BaseConf: fwConf.GetServiceConf(),
	})

	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	jobIns.Run()
}
