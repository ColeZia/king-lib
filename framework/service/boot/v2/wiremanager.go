package v2

import (
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"gl.king.im/king-lib/framework/constant"
	etcdclient "go.etcd.io/etcd/client/v3"
	goLog "log"

	"github.com/go-kratos/kratos/v2"
	"gl.king.im/king-lib/framework/config"
	"gl.king.im/king-lib/framework/service"
	"gl.king.im/king-lib/framework/service/boot"

	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
)

func NewApp(logger log.Logger, serverColl boot.NewServerCollection) *kratos.App {
	gs := serverColl.Gs
	hs := serverColl.Hs
	//appInfoInsJson, _ := json.MarshalIndent(service.AppInfoIns, "", "  ")
	appInfoInsJson, _ := json.Marshal(service.AppInfoIns)
	fmt.Printf("\nAppInfo::%s\n", appInfoInsJson)

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
	err := boot.SetTracerProvider(jaegerEp)
	if err != nil {
		//fmt.Println(err)
		panic(err)
	}
	//alerting
	boot.RegisterDefaultAlerting(serviceConf.Service.Alert)
	options := []kratos.Option{
		kratos.ID(service.AppInfoIns.Id),
		kratos.Name(service.AppInfoIns.Name),
		kratos.Version(service.AppInfoIns.Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			hs,
			gs,
		),
	}
	if serviceConf.Service.RegistryMethod == "" || serviceConf.Service.RegistryMethod == constant.Etcd { //是否将服务注册到etcd中取
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
			goLog.Fatal(err)
		}
		r := etcd.New(client)
		options = append(options, kratos.Registrar(r))
	}
	appIns := kratos.New(options...)

	return appIns
}
