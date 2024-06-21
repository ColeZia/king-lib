package config

import (
	"sync"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
)

var sourceConfSyncMap sync.Map

//获取单实例，fileNameArr实际只取第一个元素做文件名
func GetInstance(fileNameArr ...string) config.Config {

	fileName := "config"

	if len(fileNameArr) > 0 {
		fileName = fileNameArr[0]
	}

	source := "./configs/" + fileName + ".yaml"

	return GetInstanceBySource(source)
}

//获取单实例，fileNameArr实际只取第一个元素做文件名
func GetInstanceBySource(source string) config.Config {

	if val, ok := sourceConfSyncMap.Load(source); ok {
		return val.(config.Config)
	}

	instance := config.New(
		config.WithSource(
			file.NewSource(source),
		),
	)

	if err := instance.Load(); err != nil {
		panic(err)
	}

	sourceConfSyncMap.Store(source, instance)

	if val, ok := sourceConfSyncMap.Load(source); ok {
		return val.(config.Config)
	}

	return nil
}
