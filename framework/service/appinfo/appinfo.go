package appinfo

import (
	"log"
	"net/http"
	"sync"

	"gl.king.im/king-lib/framework/auth"
	"google.golang.org/grpc"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	AppInfoIns AppInfo
)

type AppInfo struct {
	Framework string
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	Flagconf string

	Id string

	Env string

	Logger       log.Logger
	HS           *http.Server
	GS           *grpc.Server
	CallScenario string
	Caller       string
	//Conf           *conf.Bootstrap
	Authentication auth.Authentication
	Authorization  auth.Authorization
	AuthRegistry   sync.Map //map[string]auth.AuthRegister
}
