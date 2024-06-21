package v1

import (
	ke "github.com/go-kratos/kratos/v2/errors"
	"gl.king.im/king-lib/framework/test/skeleton/internal/biz"

	"github.com/gin-gonic/gin"
)

type SkeletonGinService struct {
	auc *biz.SkeletonUsecase
}

func NewSkeletonGinService(auc *biz.SkeletonUsecase) *SkeletonGinService {
	return &SkeletonGinService{auc: auc}
}

func (s *SkeletonGinService) GinTest(ctx *gin.Context) (interface{}, error) {
	ctx.JSON(500, map[string]string{"ggg": "ttt"})
	return nil, ke.InternalServer("ffff", "gggg")
	return s.auc.GinTest(ctx)
}
