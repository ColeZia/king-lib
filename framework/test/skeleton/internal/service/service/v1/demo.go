package v1

import (
	"context"

	"gl.king.im/king-lib/framework/test/skeleton/internal/biz"

	//pb "boss-auth/api/auth/v1"

	pb "gl.king.im/king-lib/protobuf/api/skeleton/service/v1"
)

type SkeletonService struct {
	pb.UnimplementedSkeletonServiceServer
	auc *biz.SkeletonUsecase
}

func NewSkeletonService(auc *biz.SkeletonUsecase) *SkeletonService {
	return &SkeletonService{auc: auc}
}

func (s *SkeletonService) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetReply, error) {
	return s.auc.Get(ctx, req)
}
