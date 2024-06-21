package discovery

import (
	"sync"

	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	ke "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/registry"
	"gl.king.im/king-lib/framework/config"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func GetServiceConf() *config.Bootstrap {
	serviceConf := config.GetServiceConf()

	return serviceConf
}

var once sync.Once

var discovery registry.Discovery

func GetDiscovery() (registry.Discovery, error) {
	// c := consulAPI.DefaultConfig()
	// c.Address = conf.Consul.Address
	// c.Scheme = conf.Consul.Scheme
	// cli, err := consulAPI.NewClient(c)
	// if err != nil {
	// 	panic(err)
	// }
	// r := consul.New(cli, consul.WithHealthCheck(false))

	var onceErr error
	once.Do(func() {
		conf := GetServiceConf()

		if conf.Service.Registrys.Etcd.Addr == "" {
			onceErr = ke.InternalServer("REGISTRY_CONF_EMPTY", "系统错误")
			return
		}

		cli, err := clientv3.New(clientv3.Config{
			Endpoints: []string{conf.Service.Registrys.Etcd.Addr},
		})

		if err != nil {
			//panic(err)
			onceErr = err
		}

		discovery = etcd.New(cli)
	})

	return discovery, onceErr
}
