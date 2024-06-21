package data

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"gl.king.im/king-lib/framework/config"
	"gl.king.im/king-lib/framework/internal/di"
	tgrpc "gl.king.im/king-lib/framework/transport/grpc/client"
	admsrv "gl.king.im/king-lib/protobuf/api/admin/service/v1"
	"gl.king.im/king-lib/protobuf/api/common"
	srvpb "gl.king.im/king-lib/protobuf/api/service/service/v1"
	userv1 "gl.king.im/king-lib/protobuf/api/user/service/v1"

	//"github.com/go-kratos/kratos/pkg/naming/etcd"

	"github.com/go-redis/redis"
	"github.com/google/wire"
)

var onceInit sync.Once
var data *Data

var userSrvCli userv1.UserServiceClient
var srvSrvCli srvpb.ServiceServiceClient
var admSrvCli admsrv.AdminServiceClient

func OnceInitData() (*Data, error) {
	var initErr error
	onceInit.Do(func() {
		log.Println("OnceInitData...")
		data, initErr = NewData()
	})

	return data, initErr
}

func GetUserServiceCli() {
	return
}

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	//NewDiscovery,
	newUserServiceClient,
)

// Data .
type Data struct {
	UserCli userv1.UserServiceClient
	SrvCli  srvpb.ServiceServiceClient
}

// NewData .
func NewData() (*Data, error) {
	//conf := GetServiceConf()
	//r, err := NewDiscovery(conf.Service.Registrys)

	uc, err := newUserServiceClient()
	if err != nil {
		return nil, err
	}

	sc, err := newServiceServiceClient()
	if err != nil {
		return nil, err
	}

	return &Data{UserCli: uc, SrvCli: sc}, nil
}

func GetServiceConf() *config.Bootstrap {
	serviceConf := config.GetServiceConf()

	return serviceConf
}

//func NewDiscovery(conf *conf.Service_Registrys) (registry.Discovery, error) {
//	// c := consulAPI.DefaultConfig()
//	// c.Address = conf.Consul.Address
//	// c.Scheme = conf.Consul.Scheme
//	// cli, err := consulAPI.NewClient(c)
//	// if err != nil {
//	// 	panic(err)
//	// }
//	// r := consul.New(cli, consul.WithHealthCheck(false))
//
//	if conf.Etcd.Addr == "" {
//		return nil, ke.InternalServer("REGISTRY_CONF_EMPTY", "系统错误")
//	}
//
//	cli, err := clientv3.New(clientv3.Config{
//		Endpoints: []string{conf.Etcd.Addr},
//	})
//
//	if err != nil {
//		panic(err)
//	}
//
//	r := etcd.New(cli)
//
//	return r, nil
//}

var userSrvCliOnce sync.Once

func GetUserServiceClient() (userv1.UserServiceClient, error) {
	var onceErr error
	userSrvCliOnce.Do(func() {
		userSrvCli, onceErr = newUserServiceClient()
	})

	return userSrvCli, onceErr
}

// func newUserServiceClient(ac *conf.Auth, r registry.Discovery, tp *tracesdk.TracerProvider) userv1.UserClient {
func newUserServiceClient() (userv1.UserServiceClient, error) {
	conn, err := tgrpc.NewGrpcClientConn(common.SERVICE_NAME_BossUser.String())
	if err != nil {
		return nil, err
	}
	c := userv1.NewUserServiceClient(conn)

	return c, nil
}

var srvSrvCliOnce sync.Once

func GetServiceServiceClient() (srvpb.ServiceServiceClient, error) {
	var onceErr error
	srvSrvCliOnce.Do(func() {
		srvSrvCli, onceErr = newServiceServiceClient()
	})

	return srvSrvCli, onceErr
}

func newServiceServiceClient() (srvpb.ServiceServiceClient, error) {
	svcName := common.SERVICE_NAME_BossService.String()
	conn, err := tgrpc.NewGrpcClientConn(svcName)
	if err != nil {
		return nil, err
	}

	c := srvpb.NewServiceServiceClient(conn)

	return c, nil
}

var admSrvCliOnce sync.Once

func GetAdmServiceClient() (admsrv.AdminServiceClient, error) {
	var onceErr error
	admSrvCliOnce.Do(func() {
		admSrvCli, onceErr = newAdminServiceClient()
	})

	return admSrvCli, onceErr
}

func newAdminServiceClient() (admsrv.AdminServiceClient, error) {
	conn, err := tgrpc.NewGrpcClientConn(common.SERVICE_NAME_BossAdmin.String())
	if err != nil {
		return nil, err
	}

	stat := conn.GetState()

	fmt.Println("conn.GetState()::", stat)
	c := admsrv.NewAdminServiceClient(conn)

	return c, nil
}

var redisCliOnce sync.Once

func GetRedisClient(cnf *config.Data_Redis) (cli *redis.Client, err error) {
	redisCliOnce.Do(func() {
		cli, err = newRedisClient(cnf)
		di.SetRedisClient(cli)
	})

	cli = di.GetRedisClient()

	return
}

func newRedisClient(cnf *config.Data_Redis) (cli *redis.Client, err error) {

	if cnf == nil {
		svcCnf := config.GetServiceConf()
		if svcCnf.Data == nil || svcCnf.Data.Redis == nil {
			err = errors.New("redis未配置")
			return
		}

		cnf = svcCnf.Data.Redis
	}

	cli = redis.NewClient(&redis.Options{
		Addr:     cnf.Addr,
		Password: cnf.Password, // no password set
		DB:       0,            // use default DB
	})

	return
}
