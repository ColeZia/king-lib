package service

import (
	klog "github.com/go-kratos/kratos/v2/log"
	"gl.king.im/king-lib/framework/config"
)

type ServiceBootData struct {
	Name     string
	Version  string
	Id       string
	Logger   klog.Logger
	BaseConf *config.Bootstrap
}
