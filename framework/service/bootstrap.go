package service

import (
	"flag"
	"fmt"
	goLog "log"
	"os"
	"sync"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"gl.king.im/king-lib/framework/auth"
	"gl.king.im/king-lib/framework/config"

	etcdclient "go.etcd.io/etcd/client/v3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var appIns *kratos.App

// go build -ldflags "-X main.Version=x.y.z"
var (
	BootstrapInfoIns BootstrapInfo
	AppInfoIns       AppInfo
)

type BootstrapInfo struct {
	//	// Name is the name of the compiled software.
	//	Name string
	//	// Version is the version of the compiled software.
	//	Version string
	//	// flagconf is the config flag.
	//	Flagconf string
	//
	//	Id string
}

type AppInfo struct {
	Framework string
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	Flagconf string

	Id string

	Env string

	Logger       log.Logger
	HS           *http.Server
	GS           *grpc.Server
	CallScenario string
	Caller       string
	//Conf   conf.Bootstrap
	Conf                 *config.Bootstrap
	Authentication       auth.Authentication
	Authorization        auth.Authorization
	ScenarioAuthRegistry sync.Map //map[string]auth.AuthRegister
	AuthMethodRegistry   sync.Map //map[string]auth.AuthRegister
}

func Bootstrap(bootinfo BootstrapInfo) {
	//BootstrapInfoIns = bootinfo
	//goLog.Println("BootstrapInfoIns1111::", BootstrapInfoIns)
	//config.LoadServiceConf(AppInfoIns.Flagconf)
}

// Set global trace provider
func setTracerProvider(url string) error {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return err
	}
	tp := tracesdk.NewTracerProvider(
		// Set the sampling rate based on the parent span to 100%
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(1.0))),
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(AppInfoIns.Name),
			attribute.String("env", "dev"),
		)),
	)
	otel.SetTracerProvider(tp)
	return nil
}

func NewApp(appInfo AppInfo) *kratos.App {
	AppInfoIns = appInfo
	goLog.Println("appInfo::", AppInfoIns)
	config.LoadServiceConf(AppInfoIns.Flagconf)
	//flag.StringVar(&BootstrapInfoIns.Flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")

	flag.Parse()

	//jaeger追踪
	serviceConf := config.GetServiceConf()
	goLog.Println("serviceConf", serviceConf)
	url := serviceConf.Service.Traces.Jaeger.Endpoint
	if os.Getenv("jaeger_url") != "" {
		url = os.Getenv("jaeger_url")
	}
	if url == "" {
		goLog.Fatal("jaeger追踪器Endpoint配置为空！")
	}
	err := setTracerProvider(url)
	if err != nil {
		fmt.Println(err)
	}

	//consul服务注册
	//	client, err := api.NewClient(api.DefaultConfig())
	//	if err != nil {
	//		panic(err)
	//	}

	//etcd服务注册
	if serviceConf.Service.Registrys.Etcd.Addr == "" {
		panic("etcd服务注册Addr配置为空！")
	}
	client, err := etcdclient.New(etcdclient.Config{
		Endpoints: []string{serviceConf.Service.Registrys.Etcd.Addr},
	})

	AppInfoIns.Env = serviceConf.Service.Env

	if err != nil {
		goLog.Fatal(err)
	}

	//session.SessionStart()

	//这里的etcd是import "github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	r := etcd.New(client)

	appIns = kratos.New(
		kratos.ID(AppInfoIns.Id),
		kratos.Name(AppInfoIns.Name),
		kratos.Version(AppInfoIns.Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(AppInfoIns.Logger),
		kratos.Server(
			AppInfoIns.HS,
			AppInfoIns.GS,
		),
		//这里的consul是import consul "github.com/go-kratos/consul/registry"
		//kratos.Registrar(consul.New(client)),
		kratos.Registrar(r),
	)
	goLog.Println("appIns::", appIns.Metadata(), appIns.Endpoint(), appIns.ID(), appIns.Name())
	return appIns
}
