package server

import (
	_ "net/http/pprof"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"

	kgin "github.com/go-kratos/gin"
	"github.com/go-kratos/kratos/v2/middleware"
	mwMetadata "github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"gl.king.im/king-lib/framework/middlewares"
	fwGin "gl.king.im/king-lib/framework/middlewares/gin"
	mlog "gl.king.im/king-lib/framework/middlewares/log"
	"gl.king.im/king-lib/framework/service/boot"
)

// import "github.com/google/wire"
//
// // ProviderSet is server providers.
// var ProviderSet = wire.NewSet(NewHTTPServer, NewGRPCServer)
type Server struct {
	filters   []http.FilterFunc
	mws       []middleware.Middleware
	logger    log.Logger
	serverReg []boot.ServerRegisterCnf
	gsOpts    []grpc.ServerOption
	hsOpts    []http.ServerOption
	ginEngine *gin.Engine
}

type ServerOption func(*Server)

func WithMiddlewares(mws []middleware.Middleware) ServerOption {
	return func(o *Server) {
		o.mws = mws
	}
}

func WithFilters(httpFilters []http.FilterFunc) ServerOption {
	return func(o *Server) {
		o.filters = httpFilters
	}
}

func WithLogger(lo log.Logger) ServerOption {
	return func(o *Server) {
		o.logger = lo
	}
}

func WithGrpcServerOpts(opts []grpc.ServerOption) ServerOption {
	return func(o *Server) {
		o.gsOpts = opts
	}
}

func WithHttpServerOpts(opts []http.ServerOption) ServerOption {
	return func(o *Server) {
		o.hsOpts = opts
	}
}

func WithGinEngine(en *gin.Engine) ServerOption {
	return func(o *Server) {
		o.ginEngine = en
	}
}

func NewServer(serverReg []boot.ServerRegisterCnf, opts ...ServerOption) boot.NewServerCollection {
	server := &Server{}
	for _, v := range opts {
		v(server)
	}

	bootServerOpts := []boot.ServerOption{}
	if len(server.filters) > 0 {
		bootServerOpts = append(bootServerOpts, boot.WithServerFilters(server.filters))
	}

	if len(server.mws) > 0 {
		bootServerOpts = append(bootServerOpts, boot.WithServerMiddlewares(server.mws))
	}

	if server.ginEngine != nil {
		bootServerOpts = append(bootServerOpts, boot.WithGinEngine(server.ginEngine))
	}

	return boot.NewServer(server.logger, serverReg, server.gsOpts, server.hsOpts, bootServerOpts...)
}

func GinBaseMiddlewares() gin.HandlerFunc {
	//写在前面的先执行，后面的后执行
	return kgin.Middlewares(mwMetadata.Server(), tracing.Server(), middlewares.Recovery(), fwGin.Base, mlog.StackLog(), fwGin.HandleGinCtx)
}
