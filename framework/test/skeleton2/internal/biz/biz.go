package biz

import (
	"github.com/google/wire"
	t3Pb "gl.king.im/king-lib/framework/test/api/skeleton3/service/v1"
	"gl.king.im/king-lib/framework/transport/grpc/client"
	admPb "gl.king.im/king-lib/protobuf/api/admin/service/v1"
	"gl.king.im/king-lib/protobuf/api/common"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewSkeletonUsecase)

func NewAdminClient() (cli admPb.AdminServiceClient, err error) {
	svcName := common.SERVICE_NAME_BossAdmin.String()

	conn, err := client.NewGrpcClientConn(svcName)

	if err != nil {
		panic(err)
	}

	cli = admPb.NewAdminServiceClient(conn)

	return
}

func NewTest3Client() (cli t3Pb.Skeleton3ServiceClient, err error) {
	svcName := "BossFrameworkTest3"

	conn, err := client.NewGrpcClientConn(svcName)

	if err != nil {
		panic(err)
	}

	cli = t3Pb.NewSkeleton3ServiceClient(conn)

	return
}
