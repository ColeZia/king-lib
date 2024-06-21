package boot

import (
	"context"

	//pb "boss-auth/api/auth/v1"

	pb "gl.king.im/king-lib/protobuf/api/common/service/v1"
)

type BaseService struct {
	pb.UnimplementedBaseServiceServer
}

func NewBaseService() *BaseService {
	return &BaseService{}
}

func (s *BaseService) ListOperations(ctx context.Context, req *pb.ListOperationsRequest) (reply *pb.ListOperationsReply, err error) {
	reply = &pb.ListOperationsReply{}
	pbMap := map[string]*pb.OperationTypeList{}

	for srvFullName, opList := range operationsMap {
		pbList := &pb.OperationTypeList{}
		for _, op := range opList {
			pbList.List = append(pbList.List, &pb.OperationType{
				Operation: op.Operation,
				Summary:   op.Summary,
			})
		}
		pbMap[srvFullName] = pbList
	}
	reply.ServiceOpMap = pbMap
	return
}
