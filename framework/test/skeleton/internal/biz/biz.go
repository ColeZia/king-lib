package biz

import (
	"github.com/google/wire"
	t2Pb "gl.king.im/king-lib/framework/test/skeleton2/api/service/v1"
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

func NewTest2Client() (cli t2Pb.Skeleton2ServiceClient, err error) {
	svcName := "BossFrameworkTest2"

	conn, err := client.NewGrpcClientConn(svcName)

	if err != nil {
		panic(err)
	}

	cli = t2Pb.NewSkeleton2ServiceClient(conn)

	return
}
