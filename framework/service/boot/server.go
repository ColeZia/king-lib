package boot

import (
	"fmt"
	goHttp "net/http"
	"reflect"

	"git.e.coding.king.cloud/dev/quality/king-micro/transport/zgrpc"
	"git.e.coding.king.cloud/dev/quality/king-micro/transport/zhttp"
	"gl.king.im/king-lib/framework/constant"
	"gl.king.im/king-lib/framework/filters"
	"gl.king.im/king-lib/framework/service"

	_ "net/http/pprof"

	"github.com/gin-gonic/gin"
	kgin "github.com/go-kratos/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	mwMetadata "github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	"gl.king.im/goserver/sky-agent/v2/gather"
	v2 "gl.king.im/goserver/sky-agent/v2/gather/v2"
	skyProm "gl.king.im/goserver/sky-agent/v2/lib/client_golang/prometheus"
	"gl.king.im/king-lib/framework/coder"
	"gl.king.im/king-lib/framework/config"
	ginServer "gl.king.im/king-lib/framework/gin/server"
	"gl.king.im/king-lib/framework/interceptors"
	"gl.king.im/king-lib/framework/internal/di"
	"gl.king.im/king-lib/framework/middlewares"
	fwGin "gl.king.im/king-lib/framework/middlewares/gin"
	mlog "gl.king.im/king-lib/framework/middlewares/log"
	"gl.king.im/king-lib/framework/service/desc"
	"gl.king.im/king-lib/protobuf/api/common"
	baseSrvPb "gl.king.im/king-lib/protobuf/api/common/service/v1"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

// import "github.com/google/wire"
//
// // ProviderSet is server providers.
// var ProviderSet = wire.NewSet(NewHTTPServer, NewGRPCServer)
type HttpRegisterFunc func(s *http.Server, srv interface{})
type GrpcRegisterFunc func(s ggrpc.ServiceRegistrar, srv interface{})

type ServerRegisterCnf struct {
	ServiceImpl interface{}
	Http        interface{}                          //注册方法
	Grpc        interface{}                          //注册方法
	FileDesp    protoreflect.FileDescriptor          //文件描述
	ServiceDesp *descriptorpb.ServiceDescriptorProto //服务描述信息
	filters     []http.FilterFunc
	mw          []middleware.Middleware
}

type NewServerCollection struct {
	Gs *grpc.Server
	Hs *http.Server
}

type ServerOption func(*ServerOptions)

type ServerOptions struct {
	filters   []http.FilterFunc
	mw        []middleware.Middleware
	ginEngine *gin.Engine
}

type operationType struct {
	Operation string
	Summary   string
}

var (
	serverOption  *ServerOptions
	operationsMap = map[string][]operationType{}
)

func WithServerMiddlewares(mws []middleware.Middleware) ServerOption {
	return func(o *ServerOptions) {
		o.mw = mws
	}
}

func WithServerFilters(httpFilters []http.FilterFunc) ServerOption {
	return func(o *ServerOptions) {
		o.filters = httpFilters
	}
}

func WithGinEngine(en *gin.Engine) ServerOption {
	return func(o *ServerOptions) {
		o.ginEngine = en
	}
}

func NewServer(logger log.Logger, sc []ServerRegisterCnf, gOpts []grpc.ServerOption, hOpts []http.ServerOption, opts ...ServerOption) NewServerCollection {
	serverOption = &ServerOptions{}
	if len(opts) > 0 {
		for _, v := range opts {
			v(serverOption)
		}
	}

	//注册基础服务
	baseSrv := NewBaseService()
	sc = append(sc, ServerRegisterCnf{
		ServiceImpl: baseSrv,
		Http:        baseSrvPb.RegisterBaseServiceHTTPServer,
		Grpc:        baseSrvPb.RegisterBaseServiceServer,
		FileDesp:    baseSrvPb.File_api_common_service_v1_base_service_proto,
		ServiceDesp: protodesc.ToServiceDescriptorProto(baseSrvPb.File_api_common_service_v1_base_service_proto.Services().ByName("BaseService")),
	})

	if serverOption.ginEngine != nil {
		operations := []operationType{}
		for _, v := range serverOption.ginEngine.Routes() {
			operations = append(operations, operationType{
				Operation: v.Path,
				Summary:   v.Method + ":" + v.Handler,
			})

			operationsMap["gin"] = operations
		}
	}

	//校验
	for _, v := range sc {

		if v.ServiceImpl == nil {
			panic("服务实现不能为nil！")
		}

		st := reflect.TypeOf(v.ServiceImpl)
		sn := st.String()

		if v.Http == nil {
			panic(sn + "服务Http注册方法不能为nil！")
		}

		if v.Grpc == nil {
			panic(sn + "服务Grpc注册方法不能为nil！")
		}

		if v.ServiceDesp == nil {
			panic(sn + "服务ServiceDesp服务描述不能为nil！")
		}

		if v.FileDesp == nil {
			panic(sn + "服务FileDesp文件描述未设置！")
		}
	}

	serviceConf := config.GetServiceConf()

	if isMonitorOpen(serviceConf) {
		registerMetrics()
		serverOption.mw = append(serverOption.mw, middlewares.MonitorMiddleware(
			di.Container.RequestMetric,
			di.Container.SecondsMetric,
		))
	}

	if serviceConf.Service.InExternalNetwork {
		middlewares.SetServiceExternalNetworkState(true)
	}

	middlewares.SetUserAuthMethod(serviceConf.Service.UserAuthMethod)
	middlewares.SetSvcAuthMethod(serviceConf.Service.SvcAuthMethod)

	gs := newGRPCServer(logger, func(srv *grpc.Server) error {
		//注册grpc服务器
		fmt.Println("GRPC服务注册信息:")
		for _, v := range sc {
			//v1.RegisterAuthServer(srv, auth)
			rv := reflect.ValueOf(v.Grpc)
			args := []reflect.Value{
				reflect.ValueOf(srv),
				reflect.ValueOf(v.ServiceImpl),
			}
			rv.Call(args)

			m := map[string]*descriptorpb.MethodDescriptorProto{}

			//fileSrvsDesc := v.FileDesp.Services()
			//for i := 0; i < fileSrvsDesc.Len(); i++ {
			//	srvDesc := fileSrvsDesc.Get(i)
			//	fmt.Println("firstSrvDesc.FullName()::", srvDesc.FullName())
			//}

			operations := []operationType{}

			fileFullName := v.FileDesp.FullName()
			srvFullName := string(fileFullName) + "." + *v.ServiceDesp.Name
			fmt.Println("\nService:", srvFullName)
			for _, v2 := range v.ServiceDesp.Method {

				m[*v2.Name] = v2

				//reflectM := v.Options.ProtoReflect()
				//b, e := json.Marshal(v.Options)

				optValue := proto.GetExtension(v2.Options, common.E_BossOpts).(*common.BossOpts)

				optPrint := ""
				if optValue == nil {

				} else {
					optPrint = fmt.Sprintf("optValue: %+v", optValue)
				}

				fmt.Println("method:", *v2.Name, optPrint)

				openapiOptVal := proto.GetExtension(v2.Options, options.E_Openapiv2Operation).(*options.Operation)

				opSummary := ""
				if openapiOptVal != nil {
					opSummary = openapiOptVal.Summary
				}

				operations = append(operations, operationType{
					Operation: "/" + srvFullName + "/" + *v2.Name,
					Summary:   opSummary,
				})

				operationsMap[srvFullName] = operations

				//ServerRegCnfMap[v.ServiceDesp.]
			}

			//fmt.Println("operations::", operations)

			desc.ServerRegCnfMap[string(v.FileDesp.Package())+"."+*v.ServiceDesp.Name] = desc.ServiceDescMethodMap{
				ServDesc:  v.ServiceDesp,
				MethodMap: m,
			}

		}

		return nil
	}, gOpts)

	hs := newHTTPServer(logger, func(srv *http.Server) error {
		//注册http服务器
		for _, v := range sc {
			//v1.RegisterAuthHTTPServer(srv, auth)
			if v.Http == nil {
				panic("")
			}
			rv := reflect.ValueOf(v.Http)
			args := []reflect.Value{
				reflect.ValueOf(srv),
				reflect.ValueOf(v.ServiceImpl),
			}
			rv.Call(args)
		}
		return nil
	}, hOpts)

	if serverOption.ginEngine != nil {
		hs.HandlePrefix("/", serverOption.ginEngine)
	}

	//log.DefaultLogger.Log(log.LevelDebug, "operationsMap::", operationsMap)

	//暂不启用，已由pb定义
	if false {
		ginEngine := gin.Default()

		// 使用kratos中间件
		//ginRootRouter.Use(kgin.Middlewares(recovery.Recovery(), ginMiddleware), adminuser.BossAdminUserAuthMiddleware())
		ginEngine.Use(kgin.Middlewares(recovery.Recovery(), fwGin.AdaptToKratos, fwGin.Base))

		getOpPath := "/service/listOperations"

		ginEngine.POST(getOpPath, func(gc *gin.Context) {
			//ctx := ginServer.GinCtxToKratosCtx(gc)
			ginServer.RspSerializes(gc, map[string]interface{}{"map": operationsMap}, nil)
		})

		//hs.HandlePrefix("/", ginEngine)

		hs.Handle(getOpPath, ginEngine)
	}

	//monitor
	if isMonitorOpen(serviceConf) {
		go func() {
			addr := serviceConf.Server.Monitor.Addr
			registry := di.Container.Registry

			log.DefaultLogger.Log(log.LevelInfo, "monitor listening on:", addr)

			err := gather.RunHttpServer(addr, gather.DefaultGather, registry, nil, registry)
			if err != nil {
				panic(err)
			}
		}()
	}

	//pprof
	if serviceConf.Server.Prof.Open && serviceConf.Server.Prof.Addr != "" {
		go func() {
			log.DefaultLogger.Log(log.LevelInfo, "pprof listening on:", serviceConf.Server.Prof.Addr)
			goHttp.ListenAndServe(serviceConf.Server.Prof.Addr, nil)
		}()
	}

	//Healthy
	if serviceConf.Server.Healthy != nil && serviceConf.Server.Healthy.Open && serviceConf.Server.Healthy.Addr != "" {
		go func() {
			var healtyServer goHttp.Server
			healtyServer.Addr = serviceConf.Server.Healthy.Addr
			healtyMx := goHttp.NewServeMux()
			healtyMx.HandleFunc("/", func(w goHttp.ResponseWriter, r *goHttp.Request) {
				w.Write([]byte("ok"))
			})
			healtyServer.Handler = healtyMx
			log.DefaultLogger.Log(log.LevelInfo, "Healthy listening on:", serviceConf.Server.Healthy.Addr)
			err := healtyServer.ListenAndServe()
			if err != nil {
				panic(err)
			}

			//goHttp.ListenAndServe(serviceConf.Server.Healthy.Addr, nil)
		}()
	}

	return NewServerCollection{
		Gs: gs,
		Hs: hs,
	}

}

type GrpcServiceRgistor func(*grpc.Server) error
type HttpServiceRgistor func(*http.Server) error

func newGRPCServer(logger log.Logger, registor GrpcServiceRgistor, options []grpc.ServerOption) *grpc.Server {

	//t := reflect.TypeOf(service[0])
	//v := reflect.ValueOf(service[0])
	//log.DefaultLogger.Log(log.LevelInfo, "service reflect::", t, v)

	sc := config.GetServiceConf()
	c := sc.Server

	optionMws := []middleware.Middleware{
		//注意metadata中间件需要放在使用了metadata的中间件之前前面
		mwMetadata.Server(),
		tracing.Server(),
		//middlewares.ContextLog(),
		//recovery.Recovery(),
		middlewares.Recovery(),
		middlewares.BaseMiddleware,
		mlog.StackLog(), //需要用到baseMiddleware的用户实体，所以放在其后面
		validate.Validator(),
	}

	//if !sc.Service.CloseMetric {
	//	optionMws = append(optionMws, metrics.Server(
	//		metrics.WithSeconds(prom.NewHistogram(di.Container.SecondsMetric)),
	//		metrics.WithRequests(prom.NewCounter(di.Container.RequestMetric)),
	//	))
	//}

	optionInterceptors := []ggrpc.UnaryServerInterceptor{
		interceptors.BaseUnaryInterceptor(),
	}

	if serverOption != nil && len(serverOption.mw) > 0 {
		optionMws = append(optionMws, serverOption.mw...)
	}
	var opts []grpc.ServerOption

	svcCnf := config.GetServiceConf()
	if svcCnf.Service.RegistryMethod == constant.Istio {
		//king sdk
		zbuilder := zgrpc.NewServerBuilder()
		zbuilder = zbuilder.WithUnaryRecovery(nil)
		zOpts := zgrpc.ServerOptions(zbuilder) //king grpc opts
		//构建grpc options
		opts = []grpc.ServerOption{
			grpc.Options(zOpts...),
			grpc.UnaryInterceptor(optionInterceptors...),
			grpc.Middleware(optionMws...),
		}
	} else {
		opts = []grpc.ServerOption{
			grpc.UnaryInterceptor(optionInterceptors...),
			grpc.Middleware(optionMws...),
		}
	}

	opts = append(opts, options...)

	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}

	srv := grpc.NewServer(opts...)
	registor(srv)
	//services := srv.GetServiceInfo()

	//	for k, v := range services {
	//		log.DefaultLogger.Log(log.LevelInfo, "grpc service::", k, v.Metadata)
	//		for _, v2 := range v.Methods {
	//			log.DefaultLogger.Log(log.LevelInfo, "method::", v2)
	//		}
	//
	//		fmt.Println("")
	//
	//	}

	return srv
}

func newHTTPServer(logger log.Logger, registor HttpServiceRgistor, options []http.ServerOption) *http.Server {
	sc := config.GetServiceConf()
	c := sc.Server

	optionMws := []middleware.Middleware{
		//注意metadata中间件需要放在使用了metadata的中间件之前前面
		mwMetadata.Server(),
		tracing.Server(),
		//middlewares.ContextLog(),
		//recovery.Recovery(),
		middlewares.Recovery(),
		middlewares.BaseMiddleware,
		mlog.StackLog(),
		validate.Validator(),
		middlewares.ErrorHandler,
	}

	//if !sc.Service.CloseMetric {
	//	optionMws = append(optionMws, metrics.Server(
	//		metrics.WithSeconds(prom.NewHistogram(di.Container.SecondsMetric)),
	//		metrics.WithRequests(prom.NewCounter(di.Container.RequestMetric)),
	//	))
	//}
	var optionFilters []http.FilterFunc

	if sc.Service.DiscoveryMethod == constant.Istio {
		//king http
		httpBuilder := zhttp.NewServer(nil)
		httpBuilder = httpBuilder.WithFilters(filters.BaseFilter, filters.CorsFilter)
		kingHttpFilter := zhttp.ServerOption(httpBuilder)

		optionFilters = []http.FilterFunc{
			kingHttpFilter, //king filter
		}
	} else {
		optionFilters = []http.FilterFunc{
			filters.BaseFilter,
			filters.CorsFilter,
		}
	}

	if serverOption != nil && len(serverOption.mw) > 0 {
		optionMws = append(optionMws, serverOption.mw...)
	}

	if serverOption != nil && len(serverOption.filters) > 0 {
		optionFilters = append(optionFilters, serverOption.filters...)
	}

	var opts = []http.ServerOption{
		//http.Filter(filters.BaseFilter),
		http.Filter(optionFilters...),
		http.Middleware(optionMws...),
	}

	opts = append(opts, options...)

	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}

	opts = append(opts, http.ResponseEncoder(coder.HttpResponseEncoder()))
	//opts = append(opts, http.RequestDecoder(coder.HttpRequestDecoder()))

	srv := http.NewServer(opts...)

	//swagger文档接口
	//openAPIhandler := openapiv2.NewHandler()
	//srv.HandlePrefix("/q/", openAPIhandler)

	registor(srv)

	return srv
}

func registerMetrics() {
	requestCollector := v2.NewCounterMetricCollector("boss_request_total", "请求数", []string{"kind", "operation", "code", "reason"})
	secondsCollector := v2.NewHistogramMetricCollector("boss_request_duration_sec", "请求耗时分布", []string{"kind", "operation"},
		[]float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.250, 0.5, 1})

	di.Container.RequestMetric = requestCollector.CounterVec
	di.Container.SecondsMetric = secondsCollector.HistogramVec

	collectorList := []v2.CollectorFamily{requestCollector, secondsCollector}

	registry := skyProm.NewRegistry()
	di.Container.Registry = registry

	v2.RegisterMetricCollector(collectorList, registry, service.AppInfoIns.Name)
}

func isMonitorOpen(conf *config.Bootstrap) bool {
	return conf.Server.Monitor != nil && conf.Server.Monitor.Open
}
