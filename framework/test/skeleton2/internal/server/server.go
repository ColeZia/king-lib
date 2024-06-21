package server

import (
	"gl.king.im/king-lib/framework/scheduler"
	admv1 "gl.king.im/king-lib/framework/test/skeleton2/internal/service/admin/v1"
	srvV1 "gl.king.im/king-lib/framework/test/skeleton2/internal/service/service/v1"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
	"gl.king.im/king-lib/framework/coder"
	"gl.king.im/king-lib/framework/service/boot"
	"gl.king.im/king-lib/framework/service/boot/server"
	pb "gl.king.im/king-lib/framework/test/skeleton2/api/service/v1"
	admpb "gl.king.im/king-lib/protobuf/api/skeleton/admin/v1"
	"google.golang.org/protobuf/reflect/protodesc"
)

// ProviderSet is server providers.
// var ProviderSet = wire.NewSet(NewHTTPServer, NewGRPCServer)
var (
	ProviderSet = wire.NewSet(NewServer)
)

func NewServer(skeletonSvc *srvV1.SkeletonService, skeletonGinSvc *srvV1.SkeletonGinService, skeletonAdmSvc *admv1.SkeletonAdminService, sche scheduler.DistributedScheduler, logger log.Logger) *boot.NewServerCollection {
	var sc []boot.ServerRegisterCnf

	sc = append(sc, boot.ServerRegisterCnf{
		ServiceImpl: skeletonSvc,
		Http:        pb.RegisterSkeleton2ServiceHTTPServer,
		Grpc:        pb.RegisterSkeleton2ServiceServer,
		FileDesp:    pb.File_api_service_v1_service_proto,
		ServiceDesp: protodesc.ToServiceDescriptorProto(pb.File_api_service_v1_service_proto.Services().ByName("Skeleton2Service")),
	})

	sc = append(sc, boot.ServerRegisterCnf{
		ServiceImpl: skeletonAdmSvc,
		Http:        admpb.RegisterSkeletonAdminHTTPServer,
		Grpc:        admpb.RegisterSkeletonAdminServer,
		FileDesp:    admpb.File_api_skeleton_admin_v1_skeleton_service_proto,
		ServiceDesp: protodesc.ToServiceDescriptorProto(admpb.File_api_skeleton_admin_v1_skeleton_service_proto.Services().ByName("SkeletonAdmin")),
	})

	engine := NewEngine(skeletonGinSvc)

	serverColl := server.NewServer(sc,
		//server.WithMiddlewares([]middleware.Middleware{
		//	pkgMw.I18N(),
		//}),
		server.WithHttpServerOpts([]http.ServerOption{
			http.ErrorEncoder(coder.ErrorEncoder()),
		}),
		server.WithGinEngine(engine),
	)

	return &serverColl
}
