package boot

import (
	"flag"
	goLog "log"
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	frmCnf "gl.king.im/king-lib/framework/config"
	"gl.king.im/king-lib/framework/service"
	"gl.king.im/king-lib/framework/service/appinfo"
)

var (
	cmd string
)

type frameworkMainOpts struct {
	cmdFuncMap map[string]func(confVal interface{})
	logger     log.Logger
}

type frameworkMainOption func(*frameworkMainOpts)

func WithCmdFuncMap(cmdFuncMap map[string]func(confVal interface{})) frameworkMainOption {
	return func(o *frameworkMainOpts) {
		o.cmdFuncMap = cmdFuncMap
	}
}

func WithLogger(logger log.Logger) frameworkMainOption {
	return func(o *frameworkMainOpts) {
		o.logger = logger
	}
}

type InitAppWraper func(confVal interface{}, logger log.Logger) (*kratos.App, func(), error)

func FrameworkMain(AppInfoIns service.AppInfo, confVal interface{}, initApp InitAppWraper, opts ...frameworkMainOption) {
	options := &frameworkMainOpts{}
	for _, o := range opts {
		o(options)
	}

	service.AppInfoIns = AppInfoIns

	appinfo.AppInfoIns = appinfo.AppInfo{
		Framework: AppInfoIns.Framework,
		Name:      AppInfoIns.Name,
		Version:   AppInfoIns.Version,
	}

	flag.StringVar(&service.AppInfoIns.Flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
	flag.StringVar(&cmd, "cmd", "", "命令, eg: -cmd migrate")
	flag.Parse()
	logger := options.logger
	if logger == nil {
		logger = log.With(log.NewStdLogger(os.Stdout),
			"ts", log.DefaultTimestamp,
			"caller", log.DefaultCaller,
			"service.id", AppInfoIns.Id,
			"service.name", AppInfoIns.Name,
			"service.version", AppInfoIns.Version,
			"trace_id", tracing.TraceID(),
			"span_id", tracing.SpanID(),
		)
	}

	frmCnf.LoadServiceConf(service.AppInfoIns.Flagconf)

	serviceConf := frmCnf.GetServiceConf()
	goLog.Println("配置项检查...")
	if serviceConf.Service == nil {
		panic("服务相关配置为空！")
	}
	if serviceConf.Service.Secret == "" {
		goLog.Println("服务Secret未配置！在发起/接收接口请求时会触发panic")
	}

	if serviceConf.Service.Registrys == nil || serviceConf.Service.Registrys.Etcd == nil || serviceConf.Service.Registrys.Etcd.Addr == "" {
		goLog.Println("etcd服务注册Addr配置为空！在涉及服务发现类业务时会会触发panic，如发起接口请求时")
	}

	if serviceConf.Service.Alert == nil || serviceConf.Service.Alert.WorkWechat == nil || serviceConf.Service.Alert.WorkWechat.Hook == "" {
		goLog.Println("Service.Alert.WorkWechat.Hook配置为空！无法发送预警信息！")
	}

	_ = logger

	c := config.New(
		config.WithSource(
			file.NewSource(service.AppInfoIns.Flagconf),
		),
	)
	if err := c.Load(); err != nil {
		panic(err)
	}

	//var bc conf.Bootstrap

	if err := c.Scan(confVal); err != nil {
		panic(err)
	}

	if cmd != "" {
		if cmdFunc, ok := options.cmdFuncMap[cmd]; ok {
			cmdFunc(confVal)
			return
		} else {
			panic("此命令执行方法不存在！")
		}
	}

	//fmt.Println("service.AppInfoIns::::", service.AppInfoIns)
	app, cleanup, err := initApp(confVal, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
