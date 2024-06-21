package v1

import (
	"gl.king.im/king-lib/framework/test/skeleton3/internal/biz"

	"github.com/gin-gonic/gin"
)

type SkeletonGinService struct {
	auc *biz.SkeletonUsecase
}

func NewSkeletonGinService(auc *biz.SkeletonUsecase) *SkeletonGinService {
	return &SkeletonGinService{auc: auc}
}

func (s *SkeletonGinService) GinTest(ctx *gin.Context) (interface{}, error) {
	return s.auc.GinTest(ctx)
}
