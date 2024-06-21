package boot

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	klog "github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"

	"github.com/go-kratos/kratos/v2"
	"gl.king.im/king-lib/framework/alerting"
	"gl.king.im/king-lib/framework/config"
	"gl.king.im/king-lib/framework/internal/di"
	"gl.king.im/king-lib/framework/service"
	etcdclient "go.etcd.io/etcd/client/v3"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type fakeWriter struct {
}

func (*fakeWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

// Set global trace provider
func SetTracerProvider(url string) (err error) {

	var tExporter tracesdk.SpanExporter
	if url == "" {
		tExporter, err = stdouttrace.New(
			stdouttrace.WithWriter(&fakeWriter{}),
			// Use human-readable output.
			stdouttrace.WithPrettyPrint(),
			// Do not print timestamps for the demo.
			stdouttrace.WithoutTimestamps(),
		)
	} else {
		// Create the Jaeger exporter
		tExporter, err = jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	}

	if err != nil {
		panic("tracer provider exporter init error:" + err.Error())
	}

	tp := tracesdk.NewTracerProvider(
		// Set the sampling rate based on the parent span to 100%
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(1.0))),
		// Always be sure to batch in production.
		tracesdk.WithBatcher(tExporter),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(service.AppInfoIns.Name),
			attribute.String("env", "dev"),
		)),
	)

	otel.SetTracerProvider(tp)

	return
}

func RegisterDefaultAlerting(alert *config.Service_Alert) {
	if alert != nil {
		var channels []alerting.NotificationChannel
		if alert.WorkWechat != nil && alert.WorkWechat.Hook != "" {
			channels = append(channels, &alerting.NotiChanWorkWechat{
				Key:     "system",
				Webhook: alert.WorkWechat.Hook,
				Debug:   alert.WorkWechat.Debug,
			})
		}

		if alert.Feishu != nil && alert.Feishu.Hook != "" {
			channels = append(channels, &alerting.NotiChanFeishu{
				Key:     "system-feishu",
				Webhook: alert.Feishu.Hook,
				Debug:   alert.Feishu.Debug,
			})
		}

		if len(channels) > 0 {
			di.Container.DefaultAlerting = alerting.NewAlerting(channels, nil)
		}

	}
}

func NewApp(logger klog.Logger, serverColl NewServerCollection) *kratos.App {
	gs := serverColl.Gs
	hs := serverColl.Hs
	//appInfoInsJson, _ := json.MarshalIndent(service.AppInfoIns, "", "  ")
	appInfoInsJson, _ := json.Marshal(service.AppInfoIns)
	fmt.Printf("\nAppInfo::%s\n", appInfoInsJson)
	config.LoadServiceConf(service.AppInfoIns.Flagconf)
	//flag.StringVar(&BootstrapInfoIns.Flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")

	flag.Parse()

	//jaeger追踪
	serviceConf := config.GetServiceConf()
	//serviceConfJson, _ := json.MarshalIndent(serviceConf, "", "  ")
	serviceConfJson, _ := json.Marshal(serviceConf)
	fmt.Printf("\nServiceConf::%s\n\n", serviceConfJson)
	var jaegerEp string
	traceCnf := serviceConf.Service.Traces
	if traceCnf != nil && traceCnf.Open && traceCnf.Jaeger != nil && traceCnf.Jaeger.Endpoint != "" {
		jaegerEp = serviceConf.Service.Traces.Jaeger.Endpoint

	}

	err := SetTracerProvider(jaegerEp)

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

	if err != nil {
		log.Fatal(err)
	}

	//session.SessionStart()

	//这里的etcd是import "github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	r := etcd.New(client)

	//初始化必要组件
	//alerting
	RegisterDefaultAlerting(serviceConf.Service.Alert)

	appIns := kratos.New(
		kratos.ID(service.AppInfoIns.Id),
		kratos.Name(service.AppInfoIns.Name),
		kratos.Version(service.AppInfoIns.Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			hs,
			gs,
		),
		//这里的consul是import consul "github.com/go-kratos/consul/registry"
		//kratos.Registrar(consul.New(client)),
		kratos.Registrar(r),
	)
	return appIns
}
