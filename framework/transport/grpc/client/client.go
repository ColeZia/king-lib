package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode"

	"gl.king.im/king-lib/framework/constant"
	"gl.king.im/king-lib/framework/middlewares/auth/jwt"

	"gl.king.im/king-lib/framework/transport/grpc/discovery"

	"git.e.coding.king.cloud/dev/quality/king-micro/transport/zgrpc"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	kgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"gl.king.im/king-lib/framework/config"
	fwMd "gl.king.im/king-lib/framework/transport/metadata"
	"google.golang.org/grpc"
)

var IgnoredIP = "10.255.255.255:80"

type newGrpcClientConnOptions struct {
	Mws               []middleware.Middleware
	CliOpts           []kgrpc.ClientOption
	sdm               SvcDiscovMethod
	endpointAddr      string
	zgrpcEndpointAddr string
}

type newGrpcClientConnOption func(*newGrpcClientConnOptions)

func WithMiddleware(mws []middleware.Middleware) newGrpcClientConnOption {
	return func(opts *newGrpcClientConnOptions) {
		opts.Mws = mws
	}
}

func WithCliOpts(CliOpts []kgrpc.ClientOption) newGrpcClientConnOption {
	return func(opts *newGrpcClientConnOptions) {
		opts.CliOpts = CliOpts
	}
}

type SvcDiscovMethod string

const (
	SvcDiscovMethodEtcd          SvcDiscovMethod = "etcd"
	SvcDiscovMethodIstioDns      SvcDiscovMethod = "istio-dns"
	SvcDiscovMethodK8sDns        SvcDiscovMethod = "k8s-dns"
	SvcDiscovMethodEndpoint      SvcDiscovMethod = "endpoint"
	SvcDiscovMethodZgrpcEndpoint SvcDiscovMethod = "zgrpc-endpoint"
)

func WithSvcDiscovMethod(sdm SvcDiscovMethod) newGrpcClientConnOption {
	return func(opts *newGrpcClientConnOptions) {
		opts.sdm = sdm
	}
}

func WithDiscovEndpointAddr(epAddr string) newGrpcClientConnOption {
	return func(opts *newGrpcClientConnOptions) {
		opts.endpointAddr = epAddr
	}
}

func WithDiscovZgrpcEndpointAddr(epAddr string) newGrpcClientConnOption {
	return func(opts *newGrpcClientConnOptions) {
		opts.zgrpcEndpointAddr = epAddr
	}
}

func NewGrpcClientConn(serviceName string, opts ...newGrpcClientConnOption) (*grpc.ClientConn, error) {

	options := newGrpcClientConnOptions{}
	serviceConf := config.GetServiceConf()

	for _, o := range opts {
		o(&options)
	}

	mws := []middleware.Middleware{

		jwt.Client(), //需要改造server context的metadata，然后metadata.Client()又只执行一次，所以放在metadata中间件的前面
		metadata.Client(),
		//tracing.Client(tracing.WithTracerProvider(tp)),
		tracing.Client(),
		recovery.Recovery(),
		//jwt.Client(func(token *jwt2.Token) (interface{}, error) {
		//	return []byte(ac.ServiceKey), nil
		//}, jwt.WithSigningMethod(jwt2.SigningMethodHS256)),
	}

	if len(options.Mws) > 0 {
		mws = append(mws, options.Mws...)
	}
	var kgrpcCliOpts []kgrpc.ClientOption

	var disMethod string
	if options.sdm != "" {
		disMethod = string(options.sdm)
	} else {
		disMethod = serviceConf.Service.DiscoveryMethod
	}

	if disMethod == "" || disMethod == constant.Etcd {
		dc, err := discovery.GetDiscovery()
		if err != nil {
			return nil, err
		}
		kgrpcCliOpts = []kgrpc.ClientOption{
			kgrpc.WithEndpoint("discovery:///" + serviceName),
			kgrpc.WithDiscovery(dc),
			kgrpc.WithMiddleware(mws...),
		}
	} else if disMethod == constant.K8S { //使用k8s自身的dns
		if serviceConf.Service.ServiceAlias != nil {
			if svcName, ok := serviceConf.Service.ServiceAlias[serviceName]; ok {
				serviceName = svcName
			}
		}
		service := fmt.Sprintf("%s:9001", Camel2Case(serviceName))
		service = strings.TrimSpace(service)
		kgrpcCliOpts = []kgrpc.ClientOption{
			kgrpc.WithEndpoint(service),
			kgrpc.WithMiddleware(mws...),
		}
	} else if disMethod == string(SvcDiscovMethodEndpoint) {
		var epAddr string
		if options.endpointAddr != "" {
			epAddr = options.endpointAddr
		} else {
			if serviceConf.Service.Registrys == nil || serviceConf.Service.Registrys.Endpoint == nil || serviceConf.Service.Registrys.Endpoint.Addr == "" {
				err := errors.New("Service.Registrys.Endpoint config empty")
				return nil, err
			} else {
				epAddr = serviceConf.Service.Registrys.Endpoint.Addr
			}
		}

		kgrpcCliOpts = []kgrpc.ClientOption{
			kgrpc.WithEndpoint(epAddr),
			kgrpc.WithMiddleware(mws...),
		}
	} else if disMethod == constant.Istio { //使用istio-dns
		if serviceConf.Service.ServiceAlias != nil {
			if svcName, ok := serviceConf.Service.ServiceAlias[serviceName]; ok {
				serviceName = svcName
			}
		}
		service := fmt.Sprintf("%s.%s", Camel2Case(serviceName), serviceConf.Service.K8SNamespace)
		dialer := zgrpc.NewDialer(service)
		grpcOpts := zgrpc.ClientOptions(dialer) //构造grpcOpts
		kgrpcCliOpts = []kgrpc.ClientOption{
			kgrpc.WithEndpoint(IgnoredIP),
			kgrpc.WithOptions(grpcOpts...),
			kgrpc.WithMiddleware(mws...),
		}
	} else if disMethod == string(SvcDiscovMethodZgrpcEndpoint) {
		var epAddr string
		if options.zgrpcEndpointAddr != "" {
			epAddr = options.zgrpcEndpointAddr
		} else {
			if serviceConf.Service.Registrys == nil || serviceConf.Service.Registrys.Endpoint == nil || serviceConf.Service.Registrys.Endpoint.Addr == "" {
				err := errors.New("Service.Registrys.Endpoint config empty")
				return nil, err
			} else {
				epAddr = serviceConf.Service.Registrys.Endpoint.Addr
			}
		}

		if serviceConf.Service.ServiceAlias != nil {
			if svcName, ok := serviceConf.Service.ServiceAlias[serviceName]; ok {
				serviceName = svcName
			}
		}

		service := fmt.Sprintf("%s.%s", Camel2Case(serviceName), serviceConf.Service.K8SNamespace)
		dialer := zgrpc.NewDialer(service)
		grpcOpts := zgrpc.ClientOptions(dialer) //构造grpcOpts
		kgrpcCliOpts = []kgrpc.ClientOption{
			kgrpc.WithEndpoint(epAddr),
			kgrpc.WithOptions(grpcOpts...),
			kgrpc.WithMiddleware(mws...),
		}
	}

	if len(options.CliOpts) > 0 {
		kgrpcCliOpts = append(kgrpcCliOpts, options.CliOpts...)
	}
	conn, err := kgrpc.DialInsecure(context.Background(), kgrpcCliOpts...)

	if err != nil {
		return nil, err
	}
	return conn, nil
}

func BuildInnerMDCtx() context.Context {
	return fwMd.BuildInnerMDCtx()
}
func Camel2Case(name string) string {
	buffer := new(bytes.Buffer)
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i != 0 {
				buffer.WriteByte('-')
			}
			buffer.WriteRune(unicode.ToLower(r))
		} else {
			buffer.WriteRune(r)
		}
	}
	return buffer.String()
}
