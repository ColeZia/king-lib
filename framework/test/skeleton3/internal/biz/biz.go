package biz

import (
	"github.com/google/wire"
	"gl.king.im/king-lib/framework/transport/grpc/client"
	admPb "gl.king.im/king-lib/protobuf/api/admin/service/v1"
	"gl.king.im/king-lib/protobuf/api/common"
	payPb "gl.king.im/king-lib/protobuf/api/pay/service/v1"
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

func NewPayClient() (cli payPb.PayServiceClient, err error) {
	svcName := "BossFrameworkTest"

	conn, err := client.NewGrpcClientConn(svcName)

	if err != nil {
		panic(err)
	}

	cli = payPb.NewPayServiceClient(conn)

	return
}
