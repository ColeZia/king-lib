/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package config

import (
	"os"
	"path"

	kEtcdCnf "github.com/go-kratos/kratos/contrib/config/etcd/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"gl.king.im/king-lib/framework"
	fwConf "gl.king.im/king-lib/framework/config"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func ScanCnf(serviceName string, baseCnfPath string, val interface{}) {
	var baseCnfOpts []config.Option
	baseCnfOpts = append(baseCnfOpts, config.WithSource(
		file.NewSource(baseCnfPath),
	))

	baseCnf := config.New(baseCnfOpts...)

	defer baseCnf.Close()

	if err := baseCnf.Load(); err != nil {
		panic(err)
	}

	if err := baseCnf.Scan(fwConf.AppConfIns); err != nil {
		panic(err)
	}

	if fwConf.AppConfIns.Data != nil && fwConf.AppConfIns.Data.Config != nil {
		appCnf := fwConf.AppConfIns.Data.Config
		srcList := []config.Source{}
		for _, etcdCnfItem := range appCnf.EtcdSources {
			client, err := clientv3.New(clientv3.Config{
				Endpoints: etcdCnfItem.Endpoints,
			})

			if err != nil {
				panic(err)
			}

			//服务的配置key路径数组
			svcCnfPathArr := []string{}
			//公共配置的key路径数组
			commonCnfPathArr := []string{}

			//默认基础根目录，暂时不可配置更改
			basePath := "/service-config"

			//首先都追加基础根目录
			svcCnfPathArr = append(svcCnfPathArr, basePath)
			commonCnfPathArr = append(commonCnfPathArr, basePath)

			if appCnf.Env != fwConf.Data_Config_Env_Unknown {
				svcCnfPathArr = append(svcCnfPathArr, appCnf.Env.String())
				commonCnfPathArr = append(commonCnfPathArr, appCnf.Env.String())
			}

			//追加命名空间路径
			if appCnf.Namespace != "" {
				svcCnfPathArr = append(svcCnfPathArr, appCnf.Namespace)
				//公共配置需要判断是否跟随命名空间转移路径
				if etcdCnfItem.UseNamespaceForCommon {
					commonCnfPathArr = append(commonCnfPathArr, appCnf.Namespace)
				}
			}

			//预发布环境
			if appCnf.PreEnvDetectMethod == fwConf.Data_Config_PreEnv_OSEnv {
				OSEnvKey := "TRAFFIC_LABEL"
				if appCnf.PreOsEnvName != "" {
					OSEnvKey = appCnf.PreOsEnvName
				}

				OSEnvValue := os.Getenv(OSEnvKey)
				OSEnvValue = framework.K8sSMEnvPreName

				switch OSEnvValue {
				case framework.K8sSMEnvAlphaName:
				case framework.K8sSMEnvPreName:
					envPath := "/env-pre"
					if appCnf.PreBasePath != "" {
						envPath = appCnf.PreBasePath
					}

					svcCnfPathArr = append(svcCnfPathArr, envPath)

				case framework.K8sSMEnvProdName:

				}
			}

			svcCnfPathArr = append(svcCnfPathArr, serviceName)
			commonCnfPathArr = append(commonCnfPathArr, "Common")

			//服务的配置key路径
			svcWithPath := path.Join(svcCnfPathArr...)
			if etcdCnfItem.AbsolutePath != "" {
				svcWithPath = etcdCnfItem.AbsolutePath
			}

			//公共配置key路径
			commonCnfWithPath := path.Join(commonCnfPathArr...)
			if etcdCnfItem.CommonAbsolutePath != "" {
				commonCnfWithPath = etcdCnfItem.CommonAbsolutePath
			}

			cnfExt := ".yaml"
			switch fwConf.AppConfIns.Data.Config.Format {
			case fwConf.Data_Config_Format_Yaml:
				cnfExt = ".yaml"
			case fwConf.Data_Config_Format_Json:
				cnfExt = ".json"
			default:
				cnfExt = ".yaml"
			}

			svcWithPath += cnfExt
			commonCnfWithPath += cnfExt

			//rsp, err := client.Get(context.Background(), withPath)
			//if err != nil {
			//	panic(err)
			//}
			//fmt.Println("client.get rsp::", rsp.Kvs)

			//公共的配置源
			commonSrc, err := kEtcdCnf.New(client, kEtcdCnf.WithPath(commonCnfWithPath), kEtcdCnf.WithPrefix(true))

			if err != nil {
				panic(err)
			}

			srcList = append(srcList, commonSrc)

			//模块配置
			modulesBasePath := "modules"
			if appCnf.ModulesBasePath != "" {
				modulesBasePath = appCnf.ModulesBasePath
			}
			for _, moduleName := range appCnf.Modules {
				//模块的配置源
				modulePath := basePath + "/" + modulesBasePath + "/" + moduleName + cnfExt
				moduleSrc, err := kEtcdCnf.New(client, kEtcdCnf.WithPath(modulePath), kEtcdCnf.WithPrefix(true))
				if err != nil {
					panic(err)
				}
				srcList = append(srcList, moduleSrc)
			}

			//服务的配置源
			source, err := kEtcdCnf.New(client, kEtcdCnf.WithPath(svcWithPath), kEtcdCnf.WithPrefix(true))

			if err != nil {
				panic(err)
			}

			srcList = append(srcList, source)
		}

		if len(srcList) < 1 {
			panic("配置源列表为空")
		}

		cnf := config.New(config.WithSource(srcList...))
		defer cnf.Close()

		if err := cnf.Load(); err != nil {
			panic(err)
		}

		if err := cnf.Scan(fwConf.AppConfIns); err != nil {
			panic(err)
		}

		if err := cnf.Scan(val); err != nil {
			panic(err)
		}

	} else {
		if err := baseCnf.Scan(val); err != nil {
			panic(err)
		}
	}
}
