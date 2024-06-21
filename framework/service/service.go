package service

//此文件已废弃，请勿继续使用，new client的方法已转移至/transport目录下
import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/metadata"
	mwMetadata "github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"gl.king.im/king-lib/framework"
	"gl.king.im/king-lib/framework/auth/token"
	"gl.king.im/king-lib/framework/config"
	clientv3 "go.etcd.io/etcd/client/v3"
	goGrpc "google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type serverConnRegister struct {
	Grpc sync.Map
	Http sync.Map
}

var serverConnReg serverConnRegister

func init() {
	serverConnReg = serverConnRegister{}
}

func GenerateTokenWithConf() string {
	serviceConf := config.GetServiceConf()
	if serviceConf.Service.Secret == "" {
		panic("服务Secret未配置！")
	}

	token, err := token.GenerateToken(serviceConf.Service.Secret, AppInfoIns.Name)

	if err != nil {
		panic("token生成失败！" + err.Error())
	}

	return token
}

// etcd注册中心
var initEtcdRegOnece sync.Once

var etcdReg *etcd.Registry

func initEtcdReg() {
	initEtcdRegOnece.Do(func() {
		serviceConf := config.GetServiceConf()
		if serviceConf.Service.Registrys.Etcd.Addr == "" {
			panic("etcd服务注册Addr配置为空！")
		}
		cli, err := clientv3.New(clientv3.Config{
			Endpoints: []string{serviceConf.Service.Registrys.Etcd.Addr},
		})
		if err != nil {
			panic(err)
		}
		etcdReg = etcd.New(cli)
	})
}

func NewGrpcClientConn(serviceName string, grpcOpts ...grpc.ClientOption) (*goGrpc.ClientConn, context.Context, error) {

	callCtx := BuildMetaDataCtx()

	singleton := false
	if singleton {
		if val, ok := serverConnReg.Grpc.Load(serviceName); ok {
			connGRPC := val.(*goGrpc.ClientConn)
			if connGRPC.GetState() == connectivity.Shutdown {
				fmt.Println("connGRPC.Connect()111::", connGRPC.GetState())
				connGRPC.Connect()
			}
			fmt.Println("connGRPC.Connect()222::", connGRPC.GetState())
			return connGRPC, callCtx, nil
		}
	}

	ctx := context.Background()
	//ctx = metadata.AppendToClientContext(ctx, framework.METADATA_KEY_CALL_SCENARIO, framework.MDV_SERVICE_CALL_SCENARIO_INNER)
	initEtcdReg()

	opts := []grpc.ClientOption{
		grpc.WithEndpoint("discovery:///" + serviceName),
		grpc.WithMiddleware(
			recovery.Recovery(),
			tracing.Client(),
			mwMetadata.Client(),
		),
		grpc.WithDiscovery(etcdReg),
		grpc.WithTimeout(time.Second * 30),
	}

	if len(grpcOpts) > 0 {
		opts = append(opts, grpcOpts...)
	}

	connGRPC, err := grpc.DialInsecure(
		ctx,
		opts...,
	)

	if err != nil {
		log.Fatal(err)
	}

	if singleton {
		serverConnReg.Grpc.Store(serviceName, connGRPC)
		if connGRPC.GetState() == connectivity.Shutdown {
			fmt.Println("connGRPC.Connect()aaa::", connGRPC.GetState())
			connGRPC.Connect()
		}
		fmt.Println("connGRPC.Connect()bbb::", connGRPC.GetState())
	}
	return connGRPC, callCtx, err
}

func NewHttpClientConn(serviceName string) (*khttp.Client, context.Context, error) {
	callCtx := BuildMetaDataCtx()

	if val, ok := serverConnReg.Http.Load(serviceName); ok {
		return val.(*khttp.Client), callCtx, nil
	}

	initEtcdReg()

	ctx := context.Background()
	//ctx = metadata.AppendToClientContext(ctx, framework.METADATA_KEY_CALL_SCENARIO, framework.MDV_SERVICE_CALL_SCENARIO_INNER)

	conn, err := khttp.NewClient(
		ctx,
		khttp.WithEndpoint("discovery:///"+serviceName),
		khttp.WithDiscovery(etcdReg),
		khttp.WithTimeout(time.Second*10),
	)

	if err != nil {
		log.Fatal(err)
	}

	serverConnReg.Http.Store(serviceName, conn)

	return conn, callCtx, err
}

func BuildMetaDataCtx() context.Context {
	callCtx := context.Background()
	callCtx = metadata.AppendToClientContext(callCtx, framework.METADATA_KEY_CALL_SCENARIO, framework.MDV_SERVICE_CALL_SCENARIO_INNER)
	if AppInfoIns.Authentication == nil {
		//panic("认证器未注册实现！")
	}
	//authToken := AppInfoIns.Authentication.GetToken()
	authToken := GenerateTokenWithConf()
	callCtx = metadata.AppendToClientContext(callCtx, framework.METADATA_KEY_AUTH_TOKEN, authToken)

	return callCtx
}

func SetScenarioAuthRegistry(ar sync.Map) {
	AppInfoIns.ScenarioAuthRegistry = ar
}

func SetAuthMethodRegistry(ar sync.Map) {
	AppInfoIns.AuthMethodRegistry = ar
}
