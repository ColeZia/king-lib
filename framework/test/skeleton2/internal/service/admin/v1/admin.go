package v1

import (
	"context"

	"gl.king.im/king-lib/framework/test/skeleton2/internal/biz"

	//pb "boss-auth/api/auth/v1"

	pb "gl.king.im/king-lib/protobuf/api/skeleton/admin/v1"
)

type SkeletonAdminService struct {
	pb.UnimplementedSkeletonAdminServer
	auc *biz.SkeletonUsecase
}

func NewSkeletonAdminService(auc *biz.SkeletonUsecase) *SkeletonAdminService {
	return &SkeletonAdminService{auc: auc}
}

func (s *SkeletonAdminService) AuthGet(ctx context.Context, req *pb.GetRequest) (*pb.GetReply, error) {
	return s.auc.AdminAuthGet(ctx, req)
}
