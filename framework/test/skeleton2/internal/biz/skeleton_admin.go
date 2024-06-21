package biz

import (
	"context"

	"gl.king.im/king-lib/framework/log"
	"gl.king.im/king-lib/framework/transport/metadata"

	ke "github.com/go-kratos/kratos/v2/errors"
	pb "gl.king.im/king-lib/protobuf/api/skeleton/admin/v1"
)

//框架骨架示例repository

func (uc *SkeletonUsecase) AdminAuthGet(ctx context.Context, req *pb.GetRequest) (rep *pb.GetReply, err error) {

	//cli, err := NewPayClient()

	callCtx := metadata.BuildInnerMDCtx()
	_ = callCtx

	//_, err = cli.GetPayOrderListByOrders(ctx, &v1.GetPayOrderListByOrdersRequest{})

	//uc.log.Info("Info aaa")
	//uc.log.Infoc(ctx, "Infoc bbb")
	//uc.log.Infos(ctx, "Infos-key1", 1111, "Infos-key2", 2222, 3333, 33331111)
	//uc.log.Infow("Infow-key1", 111, "Infow-key2")
	panic("ggg")
	klh := uc.log.GetKratosLogHelper()

	klh.Info("original kratos helper:::")

	uc.log.Debuga(ctx, "a-debug", 1111, "a-key2", 2222)
	uc.log.Infoc(context.Background(), "c-info", 1111, "a-key2", 2222)
	uc.log.Warn(ctx, "warn", 1111, "key2", 2222)
	uc.log.Errora(ctx, "a-error", 1111, "a-key2", 2222)

	uc.log.Infom(ctx, "m-key1", 1111, "m-key2")
	return

	lh := uc.log.WithContext(ctx)
	lh.Info("(uc *SkeletonUsecase) Get...", "val111", "key2", "val22")
	lh.Infow("key111", "val111", "key2", "val22")

	loh, ok := log.LogHelperFromServerContext(ctx)
	if ok {
		loh.Debug("LogHelperFromServerContext log")
	}

	err = ke.InternalServer("TEST", "TEST-MSG")
	//err = ke.BadRequest("", "")

	return
}
