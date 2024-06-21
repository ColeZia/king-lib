package biz

import (
	"context"

	"gl.king.im/king-lib/framework/auth/user"
	"gl.king.im/king-lib/framework/log"
	"gl.king.im/king-lib/framework/test/skeleton/internal/conf"

	"github.com/gin-gonic/gin"
	ke "github.com/go-kratos/kratos/v2/errors"
	klog "github.com/go-kratos/kratos/v2/log"
	v1 "gl.king.im/king-lib/framework/test/skeleton2/api/service/v1"
	"gl.king.im/king-lib/protobuf/api/finance/errors"
	pb "gl.king.im/king-lib/protobuf/api/skeleton/service/v1"
)

// 框架骨架示例repository
type SkeletonRepo interface {
	Get(ctx context.Context) (val string, err error)
}

type SkeletonUsecase struct {
	repo SkeletonRepo
	log  *log.Helper
}

var BizConf *conf.Biz

func NewSkeletonUsecase(repo SkeletonRepo, logger klog.Logger, bizconf *conf.Biz) *SkeletonUsecase {
	BizConf = bizconf
	uc := &SkeletonUsecase{repo: repo, log: log.NewHelper(klog.With(logger, "module", "usecase/Skeleton"))}
	return uc
}

func (uc *SkeletonUsecase) Get(ctx context.Context, req *pb.GetRequest) (rep *pb.GetReply, err error) {

	cli, err := NewTest2Client()

	_, err = cli.Get(ctx, &v1.GetRequest{})

	uc.log.Info("Info aaa")
	uc.log.Infoc(ctx, "Infoc bbb")
	uc.log.Infos(ctx, "s-key1", 1111, "s-key2", 2222)
	uc.log.Infow("w-key1", 111)

	uc.log.Infoa(ctx, "a-key1", 1111, "a-key2", 2222)

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

func (uc *SkeletonUsecase) GinTest(gc *gin.Context) (rspData interface{}, err error) {

	err = errors.ErrorStarPicSubCurrencyNotAllowed("菜单不存在")

	panic("test....")
	err = ke.InternalServer("TEST", "TEST....")
	logCtx := log.WithStackLogContext(context.Background(), klog.LevelDebug)
	uc.log.Infoa(logCtx, "WithStackLogContext", "fff")
	uc.log.Warna(logCtx, "WithStackLogContext", "fff")
	uc.log.Errora(logCtx, "WithStackLogContext", "fff")

	log.InfoStackLogs(logCtx, uc.log.Helper)

	return
	kctx := gc.Request.Context()
	_ = kctx
	userEntity, ok := user.BossOpUserFromServerContext(gc)
	_ = userEntity
	_ = ok

	userEntity1, ok1 := user.BossOpUserFromServerContext(kctx)
	_ = userEntity1
	_ = ok1
	uc.log.Infom(kctx, "append kctx")
	uc.log.Info("(uc *SkeletonUsecase) Test...")
	panic("kkk test")
	return "", nil
}
