package di

import (
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis"
	skyProm "gl.king.im/goserver/sky-agent/v2/lib/client_golang/prometheus"
	"gl.king.im/king-lib/framework/alerting"
	clientv3 "go.etcd.io/etcd/client/v3"

	admsrv "gl.king.im/king-lib/protobuf/api/admin/service/v1"
	srvpb "gl.king.im/king-lib/protobuf/api/service/service/v1"
	userv1 "gl.king.im/king-lib/protobuf/api/user/service/v1"
)

type ContainerType struct {
	DefaultAlerting *alerting.Alerting
	SecondsMetric   *skyProm.HistogramVec
	RequestMetric   *skyProm.CounterVec
	Registry        *skyProm.Registry
	etcdClient      *clientv3.Client
	redisCli        *redis.Client

	userSrvCli userv1.UserServiceClient
	srvSrvCli  srvpb.ServiceServiceClient
	admSrvCli  admsrv.AdminServiceClient

	logger klog.Logger
}

var Container = &ContainerType{}

func GetEtcdClient() *clientv3.Client {
	return Container.etcdClient
}

func SetEtcdClient(cli *clientv3.Client) error {
	Container.etcdClient = cli
	return nil
}

func GetRedisClient() *redis.Client {
	return Container.redisCli
}

func SetRedisClient(cli *redis.Client) error {
	Container.redisCli = cli
	return nil
}

func GetLogger() klog.Logger {
	return Container.logger
}

func SetLogger(lo klog.Logger) error {
	Container.logger = lo
	return nil
}
